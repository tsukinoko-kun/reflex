package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/tsukinoko-kun/reflex/internal/config"
)

var rootCmd = &cobra.Command{
	Use:   "reflex",
	Short: "A brief description of your application",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(2)
	}
}

func getReflexConfig() (config.Config, error) {
	var conf config.Config

	wd, err := os.Getwd()
	if err != nil {
		return conf, fmt.Errorf("failed to get working directory: %v", err)
	}

	configFile, err := os.Open(filepath.Join(wd, "reflex.yaml"))
	if err != nil {
		return conf, fmt.Errorf("failed to open config file: %v", err)
	}

	if err := yaml.NewDecoder(configFile).Decode(&conf); err != nil {
		_ = configFile.Close()
		return conf, fmt.Errorf("failed to decode config file: %v", err)
	}
	_ = configFile.Close()

	return conf, nil
}
