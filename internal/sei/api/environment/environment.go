package environment

type EnvironmentType string

const (
	EnvironmentProduction EnvironmentType = "prod"
	EnvironmentTest       EnvironmentType = "test"
)

type EnvironmentConfig struct {
	Host string
}

var Environments = map[EnvironmentType]EnvironmentConfig{
	EnvironmentProduction: EnvironmentConfig{
		Host: "ksef.mf.gov.pl",
	},
	EnvironmentTest: EnvironmentConfig{
		Host: "ksef-test.mf.gov.pl",
	},
}
