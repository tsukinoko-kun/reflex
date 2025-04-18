package build

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/goccy/go-yaml"
	"github.com/tsukinoko-kun/reflex/internal/config"
	"github.com/tsukinoko-kun/reflex/internal/pacman"
)

var jsExtensions = []string{".js", ".jsx", ".ts", ".tsx"}

//go:embed frontend.nogo
var frontend_go []byte

func esbuildFile(wg *sync.WaitGroup, path string, d os.DirEntry, frontendDir string, outputDir string) {
	defer wg.Done()

	ext := filepath.Ext(d.Name())
	if !slices.Contains(jsExtensions, ext) {
		return
	}

	rel, err := filepath.Rel(frontendDir, path)
	if err != nil {
		fmt.Printf("failed to get relative path: %v\n", err)
		os.Exit(1)
		return
	}

	userImportPath := "./" + filepath.ToSlash(rel)

	wrapperContent := fmt.Sprintf(`import { createRoot } from "react-dom/client";
import { StrictMode } from "react";
import React from "react";
import App from %q;
const root = createRoot(document.getElementById('root'));
root.render(<StrictMode><App /></StrictMode>);
`, userImportPath)

	bundleFileName := strings.TrimSuffix(rel, ext) + ".js"
	outFilePath := filepath.Join(outputDir, filepath.FromSlash(bundleFileName))

	_ = os.MkdirAll(filepath.Dir(outFilePath), os.ModePerm)

	result := esbuild.Build(esbuild.BuildOptions{
		Bundle:   true,
		Write:    false,
		Platform: esbuild.PlatformBrowser,
		Format:   esbuild.FormatESModule,
		JSX:      esbuild.JSXAutomatic,
		Stdin: &esbuild.StdinOptions{
			Contents:   wrapperContent,
			ResolveDir: frontendDir,
			Sourcefile: "virtual-wrapper.tsx",
			Loader:     esbuild.LoaderTSX,
		},
		Outfile:  outFilePath,
		LogLevel: esbuild.LogLevelWarning,
		Loader: map[string]esbuild.Loader{
			".js":  esbuild.LoaderJS,
			".jsx": esbuild.LoaderJSX,
			".ts":  esbuild.LoaderTS,
			".tsx": esbuild.LoaderTSX,
			".svg": esbuild.LoaderDataURL,
		},
	})

	// Log any build errors.
	if len(result.Errors) > 0 {
		fmt.Printf("Errors bundling %s\n:", path)
		for _, e := range result.Errors {
			log.Println(e.Text)
		}
		// Continue processing other files.
		return
	}

	// Write out the first (and usually only) output file.
	if len(result.OutputFiles) > 0 {
		if err := os.WriteFile(outFilePath, result.OutputFiles[0].Contents, 0644); err != nil {
			fmt.Printf("failed to write bundle for %s: %v\n", path, err)
			os.Exit(0)
			return
		}
	}

	return
}

func Bundle() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}

	configFile, err := os.Open(filepath.Join(wd, "reflex.yaml"))
	if err != nil {
		return fmt.Errorf("failed to open config file: %v", err)
	}

	var conf config.Config
	if err := yaml.NewDecoder(configFile).Decode(&conf); err != nil {
		_ = configFile.Close()
		return fmt.Errorf("failed to decode config file: %v", err)
	}
	_ = configFile.Close()

	frontendDir := filepath.Join(wd, conf.FrontendDir, "routes")
	outputDir := filepath.Join(wd, conf.OutputDir, "frontend")
	_ = os.RemoveAll(outputDir)

	_ = os.MkdirAll(outputDir, os.ModePerm)

	wg := &sync.WaitGroup{}

	err = filepath.WalkDir(frontendDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		wg.Add(1)
		go esbuildFile(wg, path, d, frontendDir, outputDir)

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking the frontend directory: %v", err)
	}

	// write frontend.go
	if err := os.WriteFile(filepath.Join(wd, conf.OutputDir, "frontend", "frontend.go"), frontend_go, 0644); err != nil {
		fmt.Printf("failed to write frontend.go: %v\n", err)
	}

	npm, err := pacman.DetectNodePackageManager()
	if err != nil {
		return fmt.Errorf("failed to detect node package manager: %v", err)
	}

	if err := npm.ExecSilent(
		"tailwindcss",
		"-i", filepath.Join(wd, conf.FrontendDir, "style.css"),
		"-o", filepath.Join(wd, conf.OutputDir, "frontend", "style.css"),
	); err != nil {
		return fmt.Errorf("failed to build style.css with tailwindcss: %v", err)
	}

	wg.Wait()

	return nil
}
