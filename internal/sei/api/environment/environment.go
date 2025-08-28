package environment

import (
	"ksef/internal/config"
	"path"
)

type EnvironmentType string

const (
	EnvironmentProduction EnvironmentType = "prod"
	EnvironmentTest       EnvironmentType = "test"
)

type EnvironmentConfig struct {
	Host            string
	CertificateFile string
}

var Environments = map[EnvironmentType]EnvironmentConfig{
	EnvironmentProduction: EnvironmentConfig{
		Host: "ksef.mf.gov.pl",
	},
	EnvironmentTest: EnvironmentConfig{
		Host: "ksef-test.mf.gov.pl",
	},
}

func (ec EnvironmentConfig) GetCertificateFile(cfg config.Config) string {
	return path.Join(cfg.CertificatesPath, ec.Host+".der")
}
