package runtime

import (
	"fmt"
	"ksef/core/ksef"
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

// DEPRECATED - definetely remove in v0.3
var LegacyEnvironmentHosts = map[string][]string{
	TestEnvironmentId: {"ksef-test.mf.gov.pl", "api-test.ksef.mf.gov.pl"},
	DemoEnvironmentId: {"ksef-demo.mf.gov.pl", "api-demo.ksef.mf.gov.pl"},
	ProdEnvironmentId: {"ksef.mf.gov.pl"},
}

var environments = map[string]Environment{
	TestEnvironmentId: testEnvironment,
	DemoEnvironmentId: demoEnvironment,
	ProdEnvironmentId: prodEnvironment,
}

func Gateway(vip *viper.Viper) (ksef.KSeFGateway, error) {
	// TODO: get rid of the GetEnvironmentId
	return ksef.ParseGateway(GetEnvironmentId(vip))
}

func GatewayConfig(vip *viper.Viper) (ksef.KSeFGatewayConfig, error) {
	return ksef.GatewayConfig(GetEnvironmentId(vip))
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
