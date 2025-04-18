package config

type Config struct {
	FrontendDir string `yaml:"frontend"`
	BackendDir  string `yaml:"backend"`
	OutputDir   string `yaml:"output"`
}
