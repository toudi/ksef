package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigType struct {
	Logging map[string]string `yaml:"logging"`
}

var Config ConfigType

func ReadConfig(configFile string) error {
	file, err := os.Open(configFile)
	if err != nil {
		return fmt.Errorf("unable to open config file: %v", err)
	}
	if err = yaml.NewDecoder(file).Decode(&Config); err != nil {
		return fmt.Errorf("unable to parse config file: %v", err)
	}

	return nil
}
