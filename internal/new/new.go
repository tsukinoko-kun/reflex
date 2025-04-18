package new

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/goccy/go-yaml"
	"github.com/tsukinoko-kun/reflex/internal/config"
	"github.com/tsukinoko-kun/reflex/internal/static"
)

const (
	npm uint8 = iota
	pnpm
	yarn
)

type nodePackage struct {
	Name            string            `json:"name"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

var (
	packageNameRegex     = regexp.MustCompile(`^[a-z][a-z0-9]*$`)
	packageNameForbidden = []string{
		"reflex",
		"main",
		"init",
	}
)

func New() error {
	title := "example"
	pacman := npm
	srcDir := true
	location, err := os.Getwd()

	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	if err := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Project name").
			Validate(func(s string) error {
				if !packageNameRegex.MatchString(s) {
					return errors.New("invalid package name, must start with a lowercase letter and contain only lowercase letters and numbers")
				}
				if slices.Contains(packageNameForbidden, s) {
					return fmt.Errorf("invalid package name, %s is forbidden", s)
				}
				return nil
			}).
			Value(&title),
		huh.NewFilePicker().
			ShowHidden(true).
			ShowPermissions(true).
			CurrentDirectory(location).
			DirAllowed(true).
			FileAllowed(false).
			Title("Location").
			Value(&location),
		huh.NewSelect[uint8]().
			Title("Node package manager").
			Value(&pacman).
			Options(
				huh.NewOption("npm", npm),
				huh.NewOption("pnpm", pnpm),
				huh.NewOption("yarn", yarn),
			),
		huh.NewConfirm().
			Title("Use src directory").
			Value(&srcDir),
	)).Run(); err != nil {
		return nil
	}

	if err := os.Chdir(location); err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}

	conf := config.Config{
		OutputDir: "out",
	}
	if srcDir {
		conf.FrontendDir = "src/frontend"
		conf.BackendDir = "src/backend"
	} else {
		conf.FrontendDir = "frontend"
		conf.BackendDir = "backend"
	}

	if err := writeInstinctConfig(conf); err != nil {
		return fmt.Errorf("failed to write instinct config: %w", err)
	}

	if err := writePackageJson(nodePackage{
		Name:    "reflex",
		Scripts: map[string]string{"dev": "reflex dev", "build": "reflex build", "lint": "reflex lint"},
		Dependencies: map[string]string{
			"react":     "^19.1.0",
			"react-dom": "^19.1.0",
		},
		DevDependencies: map[string]string{
			"@types/react":     "^19.1.2",
			"@types/react-dom": "^19.1.2",
			"@biomejs/biome":   "1.9.4",
		},
	}); err != nil {
		return fmt.Errorf("failed to write package.json: %w", err)
	}

	readme := strings.Builder{}
	readme.WriteString("# ")
	readme.WriteString(title)
	readme.WriteString("\n\n## Install dependencies\n\n[Reflex](https://github.com/tsukinoko-kun/reflex)\n\n```shell\n")
	switch pacman {
	case npm:
		readme.WriteString("npm install")
	case yarn:
		readme.WriteString("yarn install")
	case pnpm:
		readme.WriteString("pnpm install")
	}
	readme.WriteString("\n```\n")

	if err := writeStringToFile("README.md", readme.String()); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}

	gitignore := strings.Builder{}
	gitignore.WriteString(".DS_Store\n")
	gitignore.WriteString(".idea/\n")
	gitignore.WriteString(".vscode/\n")
	gitignore.WriteString("node_modules/\n")
	gitignore.WriteString(conf.OutputDir + "/\n")
	switch pacman {
	case npm:
		gitignore.WriteString("pnpm-lock.yaml\n")
		gitignore.WriteString("pnpm-workspace.yaml\n")
		gitignore.WriteString("yarn.lock\n")
	case yarn:
		gitignore.WriteString("pnpm-lock.yaml\n")
		gitignore.WriteString("pnpm-workspace.yaml\n")
		gitignore.WriteString("package-lock.json\n")
	case pnpm:
		gitignore.WriteString("yarn.lock\n")
		gitignore.WriteString("package-lock.json\n")
	}

	if err := writeStringToFile(".gitignore", gitignore.String()); err != nil {
		return fmt.Errorf("failed to write .gitignore: %w", err)
	}

	if err := static.WriteContentTo(location); err != nil {
		return fmt.Errorf("failed to write static files: %w", err)
	}

	if err := static.WriteFrontendTo(filepath.Join(location, conf.FrontendDir)); err != nil {
		return fmt.Errorf("failed to write frontend files: %w", err)
	}

	if err := static.WriteBackendTo(filepath.Join(location, conf.BackendDir)); err != nil {
		return fmt.Errorf("failed to write backend files: %w", err)
	}

	var cmd *exec.Cmd
	switch pacman {
	case npm:
		cmd = exec.Command("npm", "install")
	case yarn:
		cmd = exec.Command("yarn", "install")
	case pnpm:
		cmd = exec.Command("pnpm", "install")
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run command: %w", err)
	}

	cmd = exec.Command("go", "mod", "init", title)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run command: %w", err)
	}
	cmd = exec.Command("go", "mod", "tidy", title)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run command: %w", err)
	}

	return nil
}

func writeStringToFile(filename string, content string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("failed to write string to file: %w", err)
	}

	return nil
}

func writeInstinctConfig(conf config.Config) error {
	f, err := os.Create("instinct.yaml")
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer f.Close()

	if err := yaml.NewEncoder(f).Encode(&conf); err != nil {
		return fmt.Errorf("failed to encode config file: %w", err)
	}

	return nil
}

func writePackageJson(conf nodePackage) error {
	f, err := os.Create("package.json")
	if err != nil {
		return fmt.Errorf("failed to create package.json file: %w", err)
	}
	defer f.Close()

	je := json.NewEncoder(f)
	je.SetIndent("", "  ")
	if err := je.Encode(&conf); err != nil {
		_ = f.Close()
		return fmt.Errorf("failed to encode package.json file: %w", err)
	}

	return nil
}
