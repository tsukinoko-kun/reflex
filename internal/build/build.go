package build

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/goccy/go-yaml"
	"github.com/tsukinoko-kun/reflex/internal/config"
)

var jsExtensions = []string{".js", ".jsx", ".ts", ".tsx"}

//go:embed frontend.nogo
var frontend_go []byte

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

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	err = filepath.WalkDir(frontendDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(d.Name())
		if !slices.Contains(jsExtensions, ext) {
			return nil
		}

		rel, err := filepath.Rel(frontendDir, path)
		if err != nil {
			return err
		}

		userImportPath := "./" + filepath.ToSlash(rel)

		wrapperContent := `import React from 'react';
import ReactDOM from 'react-dom';
import Component from '` + userImportPath + `';

ReactDOM.render(
  <Component />,
  document.getElementById('root')
);
`

		bundleFileName := strings.TrimSuffix(rel, ext) + ".js"
		outFilePath := filepath.Join(outputDir, filepath.FromSlash(bundleFileName))

		if err := os.MkdirAll(filepath.Dir(outFilePath), os.ModePerm); err != nil {
			return err
		}

		result := esbuild.Build(esbuild.BuildOptions{
			Bundle: true,
			Write:  false,
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
			return nil
		}

		// Write out the first (and usually only) output file.
		if len(result.OutputFiles) > 0 {
			if err := os.WriteFile(outFilePath, result.OutputFiles[0].Contents, 0644); err != nil {
				fmt.Printf("failed to write bundle for %s: %v", path, err)
			}
		}

		return nil
	})

	// write frontend.go
	if err := os.WriteFile(filepath.Join(wd, conf.OutputDir, "frontend", "frontend.go"), frontend_go, 0644); err != nil {
		fmt.Printf("failed to write frontend.go: %v", err)
	}

	if err != nil {
		return fmt.Errorf("error walking the frontend directory: %v", err)
	}
	return nil
}
