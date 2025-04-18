package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/tsukinoko-kun/reflex/internal/build"
)

var fullBuildMut = sync.Mutex{}

func fullBuild(trackTime bool) {
	fullBuildMut.Lock()
	defer fullBuildMut.Unlock()

	startTime := time.Now()

	if err := build.Bundle(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reflexConf, err := getReflexConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	goBuildCmd := exec.Command(
		"go",
		"build",
		"-o", "./"+filepath.Join(reflexConf.OutputDir, "backend", binName),
		"./"+reflexConf.BackendDir,
	)
	goBuildCmd.Stderr = os.Stderr
	goBuildCmd.Stdout = os.Stdout
	if err := goBuildCmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if trackTime {
		endTime := time.Now()
		fmt.Printf("Build completed in %v\n", endTime.Sub(startTime))
	}
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {
		fullBuild(true)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
