package runtime

import (
	"fmt"
	"ksef/internal/utils"

	"github.com/spf13/viper"
)

type Environment struct {
	ID           string
	API          string // bramka API
	QRCode       string // url do weryfikacji kodów QR
	NIPValidator utils.NIPValidatorType
}

const (
	cfgKeyEnvironment = "environment"

	ProdEnvironmentId = "prod"
	DemoEnvironmentId = "demo"
	TestEnvironmentId = "test"
)

var (
	testEnvironment = Environment{
		ID:           TestEnvironmentId,
		API:          "https://api-test.ksef.mf.gov.pl",
		QRCode:       "https://qr-test.ksef.mf.gov.pl",
		NIPValidator: utils.NIPLengthValidator,
	}
	demoEnvironment = Environment{
		ID:           DemoEnvironmentId,
		API:          "https://api-demo.ksef.mf.gov.pl",
		QRCode:       "https://qr-demo.ksef.mf.gov.pl",
		NIPValidator: utils.NIPValidator,
	}
	prodEnvironment = Environment{
		ID:           ProdEnvironmentId,
		API:          "https://api.ksef.mf.gov.pl",
		QRCode:       "https://qr.ksef.mf.gov.pl",
		NIPValidator: utils.NIPValidator,
	}
)

var environments = map[string]Environment{
	TestEnvironmentId: testEnvironment,
	DemoEnvironmentId: demoEnvironment,
	ProdEnvironmentId: prodEnvironment,
}

var LegacyEnvironmentHosts = map[string][]string{
	TestEnvironmentId: {"ksef-test.mf.gov.pl", "api-test.ksef.mf.gov.pl"},
	DemoEnvironmentId: {"ksef-demo.mf.gov.pl", "api-demo.ksef.mf.gov.pl"},
	ProdEnvironmentId: {"ksef.mf.gov.pl"},
}

func GetEnvironmentId(vip *viper.Viper) string {
	return vip.GetString(cfgKeyEnvironment)
}

func GetEnvironment(vip *viper.Viper) Environment {
	return environments[GetEnvironmentId(vip)]
}

func SetEnvironment(vip *viper.Viper, env string) {
	vip.Set(cfgKeyEnvironment, env)
}

func GetNIPValidator(vip *viper.Viper) (utils.NIPValidatorType, error) {
	environmentId := vip.GetString(cfgKeyEnvironment)
	if environment, exists := environments[environmentId]; exists {
		return environment.NIPValidator, nil
	}

	return nil, fmt.Errorf("unknown environment: %s", environmentId)
}
