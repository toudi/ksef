package runtime

import (
	"ksef/cmd/ksef/flags"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var cfgKeyNip = "nip"

func GetNIP(vip *viper.Viper) (string, error) {
	nipValidator, err := GetNIPValidator(vip)
	if err != nil {
		return "", err
	}
	rawNIP := viper.GetString(cfgKeyNip)
	return rawNIP, nipValidator(rawNIP)
}

func SetNIP(vip *viper.Viper, nip string) {
	vip.Set(cfgKeyNip, nip)
}

func SetNIPFromFlags(vip *viper.Viper, _flags *pflag.FlagSet) error {
	nipValue, err := _flags.GetString(flags.FlagNameNIP)
	if err != nil {
		return err
	}
	SetNIP(vip, nipValue)
	return nil
}
