package runtime

import "github.com/spf13/viper"

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
