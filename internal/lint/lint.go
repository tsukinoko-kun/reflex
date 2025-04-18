package lint

import (
	"fmt"
	"os"
	"os/exec"

	"slices"

	"github.com/tsukinoko-kun/reflex/internal/pacman"
)

var (
	biomeConfigFileNames  = []string{"biome.json", "biome.jsonc"}
	eslintConfigFileNames = []string{"eslint.config.js", "eslint.config.mjs", "eslint.config.cjs", "eslint.config.ts", "eslint.config.mts", "eslint.config.cts"}
)

func exists(dirEntries []os.DirEntry, names []string) bool {
	for _, entry := range dirEntries {
		if slices.Contains(names, entry.Name()) {
			return true
		}
	}
	return false
}

func Format() error {
	npm, err := pacman.DetectNodePackageManager()
	if err != nil {
		return err
	}

	_ = npm.ExecSilent("biome", "check", "--write")
	cmd := exec.Command("go", "fmt", "./...")
	_ = cmd.Run()

	return nil
}

func Run(fix bool) error {
	npm, err := pacman.DetectNodePackageManager()
	if err != nil {
		return err
	}

	dirEntries, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read directory entries: %w", err)
	}

	// is biome configured?
	if exists(dirEntries, biomeConfigFileNames) {
		if fix {
			if err := npm.Exec("biome", "check", "--write", "--unsafe"); err != nil {
				return fmt.Errorf("failed to lint with biome: %w", err)
			}
		} else {
			if err := npm.Exec("biome", "check"); err != nil {
				return fmt.Errorf("failed to lint with biome: %w", err)
			}
		}
	}

	// is eslint configured?
	if exists(dirEntries, eslintConfigFileNames) {
		if fix {
			if err := npm.Exec("eslint", "--fix"); err != nil {
				return fmt.Errorf("failed to lint with eslint: %w", err)
			}
		} else {
			if err := npm.ExecSilent("eslint"); err != nil {
				return fmt.Errorf("failed to lint with eslint: %w", err)
			}
		}
	}

	return nil
}
