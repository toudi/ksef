package runtime

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	cfgKeyIgnoreSSLErrors = "https.ignore-ssl-errors"
)

func FlagIgnoreSSLErrors(flagSet *pflag.FlagSet) {
	flagSet.Bool(cfgKeyIgnoreSSLErrors, false, "ignoruj błędy SSL (błąd uwierzytelnienia certyfikatu)")
}

func IgnoreSSLErrors(vip *viper.Viper) bool {
	return vip.GetBool(cfgKeyIgnoreSSLErrors)
}
