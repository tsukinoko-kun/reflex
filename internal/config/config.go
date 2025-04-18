package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type Config struct {
	FrontendDir string `yaml:"frontend"`
	BackendDir  string `yaml:"backend"`
	PublicDir   string `yaml:"public"`
	OutputDir   string `yaml:"output"`
}

var (
	WD, _ = os.Getwd()
	Conf  *Config
)

func Load() *Config {
	if Conf != nil {
		return Conf
	}

	WD, _ = os.Getwd()
	var cfg Config
	f, err := os.Open(filepath.Join(WD, "reflex.yaml"))
	if err != nil {
		fmt.Printf("error loading config: %v\n", err)
		os.Exit(1)
		return nil
	}
	defer f.Close()
	yd := yaml.NewDecoder(f)
	if err := yd.Decode(&cfg); err != nil {
		fmt.Printf("error decoding config: %v\n", err)
		os.Exit(1)
		return nil
	}

	Conf = &cfg
	return Conf
}
