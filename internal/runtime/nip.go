package runtime

import (
	"errors"
	"ksef/cmd/ksef/flags"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfgKeyNip        = "nip"
	errValidatingNIP = errors.New("błąd walidacji numeru NIP")
)

func GetNIP(vip *viper.Viper) (string, error) {
	nipValidator, err := GetNIPValidator(vip)
	if err != nil {
		return "", err
	}
	rawNIP := viper.GetString(cfgKeyNip)
	return rawNIP, nipValidator(rawNIP)
}

func SetNIPIfUnset(vip *viper.Viper, nip string) {
	if !vip.IsSet(cfgKeyNip) {
		SetNIP(vip, nip)
	}
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

func NIPSlice(flagSet *pflag.FlagSet) {
	flagSet.StringSliceP(cfgKeyNip, "n", nil, "numery NIP podmiotów")
}

func GetNIPSlice(vip *viper.Viper) ([]string, error) {
	var nipSlice []string
	nipValidator, err := GetNIPValidator(vip)
	if err != nil {
		return nil, err
	}
	for _, nip := range vip.GetStringSlice(cfgKeyNip) {
		if err := nipValidator(nip); err != nil {
			return nil, errors.Join(errValidatingNIP, errors.New(nip))
		}
		nipSlice = append(nipSlice, nip)
	}

	return nipSlice, nil
}
