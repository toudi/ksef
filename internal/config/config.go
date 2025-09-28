package config

import (
	"fmt"
	"ksef/internal/certsdb"
	"ksef/internal/environment"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type APIConfig struct {
	Environment    environment.Config
	CertificatesDB *certsdb.CertificatesDB
}

type Config struct {
	Logging     map[string]string `yaml:"logging"`
	PDFRenderer map[string]string `yaml:"pdf-renderer"`
}

func (c Config) APIConfig(env environment.Environment) APIConfig {
	certsDB, err := certsdb.OpenOrCreate(env)
	if err != nil {
		log.Fatalf("unable to open certificates db: %v", err)
	}
	return APIConfig{
		Environment:    environment.GetConfig(env),
		CertificatesDB: certsDB,
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
