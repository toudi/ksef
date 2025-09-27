package config

import (
	"fmt"
	"ksef/internal/utils"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type APIEnvironment string

const (
	APIEnvironmentTest APIEnvironment = "ksef-test.mf.gov.pl"
	APIEnvironmentProd APIEnvironment = "ksef.mf.gov.pl"
)

var nipValidators = map[APIEnvironment]utils.NIPValidatorType{
	APIEnvironmentProd: utils.NIPValidator,
	APIEnvironmentTest: utils.NIPLengthValidator,
}

type Certificate string

func (c Certificate) DER() string {
	return string(c) + ".der"
}

func (c Certificate) PEM() string {
	return string(c) + ".pem"
}

type APIConfig struct {
	Host         string
	Certificate  Certificate
	NIPValidator utils.NIPValidatorType
}

type Config struct {
	Logging          map[string]string `yaml:"logging"`
	PDFRenderer      map[string]string `yaml:"pdf-renderer"`
	CertificatesPath string            `yaml:"certificates-path"`
}

func (c Config) APIConfig(env APIEnvironment) APIConfig {
	return APIConfig{
		Host:         string(env),
		Certificate:  Certificate(path.Join(c.CertificatesPath, string(env))),
		NIPValidator: nipValidators[env],
	}
}

var config = Config{
	CertificatesPath: "klucze/",
}

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
