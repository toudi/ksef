package config

import (
	"fmt"
	"ksef/internal/environment"
	"os"

	"gopkg.in/yaml.v3"
)

type APIConfig struct {
	Environment environment.Config
}

type Config struct {
	Logging     map[string]string `yaml:"logging"`
	PDFRenderer map[string]string `yaml:"pdf-renderer"`
}

func (c Config) APIConfig(env environment.Environment) APIConfig {
	return APIConfig{
		Environment: environment.GetConfig(env),
	}
}

var config = Config{}

func ReadConfig(configFile string) error {
	file, err := os.Open(configFile)
	if err != nil {
		return fmt.Errorf("unable to open config file: %v", err)
	}
	if err = yaml.NewDecoder(file).Decode(&config); err != nil {
		return fmt.Errorf("unable to parse config file: %v", err)
	}

	return nil
}

func GetConfig() Config {
	return config
}
