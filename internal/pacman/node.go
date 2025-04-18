package pacman

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type NodePacman uint8

var (
	ErrUnknownPackageManager = errors.New("unknown package manager")
)

const (
	NPM NodePacman = iota + 1
	PNPM
	YARN
	BUN
)

func (n NodePacman) String() string {
	switch n {
	case NPM:
		return "npm"
	case PNPM:
		return "pnpm"
	case YARN:
		return "yarn"
	case BUN:
		return "bun"
	default:
		return "unknown"
	}
}

func DetectNodePackageManager() (NodePacman, error) {
	// use package lock files to detect package manager
	if _, err := os.Stat("pnpm-lock.yaml"); err == nil {
		return PNPM, nil
	}
	if _, err := os.Stat("yarn.lock"); err == nil {
		return YARN, nil
	}
	if _, err := os.Stat("bun.lockb"); err == nil {
		return BUN, nil
	}
	if _, err := os.Stat("package-lock.json"); err == nil {
		return NPM, nil
	}
	return 0, ErrUnknownPackageManager
}

func (n NodePacman) InstallDependencies() error {
	cmd := exec.Command(n.String(), "install")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install dependencies via %s: %w", n.String(), err)
	}
	return nil
}

func (n NodePacman) Run(script string) error {
	cmd := exec.Command(n.String(), "run", script)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run script via %s: %w", n.String(), err)
	}
	return nil
}

func (n NodePacman) ExecSilent(binary string, args ...string) error {
	cmd := exec.Command(filepath.Join("./node_modules", ".bin", binary), args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute binary via node: %w", err)
	}
	return nil
}

func (n NodePacman) Exec(binary string, args ...string) error {
	cmd := exec.Command(filepath.Join("./node_modules", ".bin", binary), args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute binary via node: %w", err)
	}
	return nil
}
