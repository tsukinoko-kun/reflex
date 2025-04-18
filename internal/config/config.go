package config

type Config struct {
	FrontendDir string `yaml:"frontend"`
	BackendDir  string `yaml:"backend"`
	PublicDir   string `yaml:"public"`
	OutputDir   string `yaml:"output"`
}
