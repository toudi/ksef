package flags

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	FlagNameNIP = "nip"
)

// NIP registers a nip parameter that binds to viper instance
func NIP(flagSet *pflag.FlagSet) {
	flagSet.FuncP(FlagNameNIP, "n", "numer NIP podmiotu", func(value string) error {
		if value != "" {
			vip := viper.GetViper()
			vip.Set(FlagNameNIP, value)
		}
		return nil
	})
}
