package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/tsukinoko-kun/reflex/internal/config"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fullBuild(true)
		devRun()

		conf := config.Load()

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Printf("failed to create fs watcher: %v\n", err)
			os.Exit(1)
		}
		defer watcher.Close()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		watchAllDirs := func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				_ = watcher.Add(path)
			}
			return nil
		}

		go devWatcher(watcher, watchAllDirs)

		filepath.WalkDir(filepath.Join(config.WD, conf.PublicDir), watchAllDirs)
		filepath.WalkDir(filepath.Join(config.WD, conf.BackendDir), watchAllDirs)
		filepath.WalkDir(filepath.Join(config.WD, conf.FrontendDir), watchAllDirs)

		<-sigs

		_ = watcher.Close()

		if devRunCmd != nil {
			_ = devRunCmd.Process.Signal(syscall.SIGHUP)
		}
	},
}

var (
	debouncerMut   = sync.Mutex{}
	debouncerTimer *time.Timer
)

func debouncedFullBuild() {
	debouncerMut.Lock()
	defer debouncerMut.Unlock()

	if debouncerTimer != nil {
		debouncerTimer.Stop()
	}

	debouncerTimer = time.NewTimer(500 * time.Millisecond)
	select {
	case <-debouncerTimer.C:
		fullBuild(true)
		devRun()
	}
}

func devWatcher(watcher *fsnotify.Watcher, watchAllDirs func(string, fs.DirEntry, error) error) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			fileInfo, err := os.Stat(event.Name)
			if err == nil && fileInfo.IsDir() {
				filepath.WalkDir(event.Name, watchAllDirs)
			}
			if event.Has(fsnotify.Write | fsnotify.Create | fsnotify.Remove | fsnotify.Rename) {
				debouncedFullBuild()
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("error:", err)
		}
	}
}

var devRunCmd *exec.Cmd

func devRun() {
	conf := config.Load()

	if devRunCmd != nil {
		devRunCmd.Process.Signal(syscall.SIGTERM)
		devRunCmd.Wait()
		devRunCmd = nil
	}

	devRunCmd = exec.Command(filepath.Join(config.WD, conf.OutputDir, "backend", binName), "--addr", ":4321")
	devRunCmd.Stdout = os.Stdout
	devRunCmd.Stderr = os.Stderr
	devRunCmd.Stdin = os.Stdin
	if err := devRunCmd.Start(); err != nil {
		fmt.Printf("failed to start dev server: %v\n", err)

	}
}

func init() {
	rootCmd.AddCommand(devCmd)
}
