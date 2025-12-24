package runtime

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	cfgKeyPreferredCertProfile = "cert-profile"
)

func CertProfileFlag(flagSet *pflag.FlagSet) {
	flagSet.String(cfgKeyPreferredCertProfile, "", "ID preferowanego profilu certyfikatów w ramach podmiotu")
}

func GetPreferredCertProfile(vip *viper.Viper) string {
	return vip.GetString(cfgKeyPreferredCertProfile)
}
