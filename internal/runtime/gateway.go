package runtime

import (
	"ksef/internal/utils"

	"github.com/spf13/viper"
)

type Gateway string

const (
	cfgKeyGateway = "gateway"

	ProdGateway Gateway = "ksef.mf.gov.pl"
	TestGateway Gateway = "ksef-test.mf.gov.pl"
	DemoGateway Gateway = "ksef-demo.mf.gov.pl"
)

func GetGateway(vip *viper.Viper) Gateway {
	return Gateway(vip.GetString(cfgKeyGateway))
}

func SetGateway(vip *viper.Viper, gw Gateway) {
	vip.Set(cfgKeyGateway, gw)
}

func GetNIPValidator(vip *viper.Viper) (utils.NIPValidatorType, error) {
	var validator = utils.NIPLengthValidator

	gateway := GetGateway(vip)

	if gateway != TestGateway {
		validator = utils.NIPValidator
	}

	return validator, nil
}
