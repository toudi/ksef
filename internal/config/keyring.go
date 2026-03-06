package config

import (
	"ksef/internal/flags"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(flags.CfgKeyKeyringEngine, flags.KeyringEngineSystem)
}

func FileKeyringFlags(flagSet *pflag.FlagSet) {
	flagSet.String(flags.CfgKeyKeyringFileLocation, "", "ścieżka do keyringu opartego o plik")
	flagSet.Bool(flags.CfgKeyKeyringFileAskPassword, false, "pytaj o hasło do keyringu na wejściu standardowym (stdin)")
	flagSet.Bool(flags.CfgKeyKeyringFileBuffered, false, "buforuj keyring w pamięci")
	flagSet.String(flags.CfgKeyKeyringFilePasswordFile, "", "ścieżka do pliku z hasłem keyringu")
	flagSet.String(flags.CfgKeyKeyringFilePasswordEnvVar, "", "nazwa zmiennej środowiskowej która zawiera hasło do keyringu")
}

func KeyringFlags(flagSet *pflag.FlagSet) {
	flagSet.String(flags.CfgKeyKeyringEngine, "system", "silnik keyringu")
	FileKeyringFlags(flagSet)
}
