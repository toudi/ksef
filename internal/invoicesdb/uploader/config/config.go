package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	cfgKeyUseBatchSession = "upload.batch"
)

type UploaderConfig struct {
	BatchSession bool
}

func UploaderFlags(flagSet *pflag.FlagSet) {
	flagSet.Bool(cfgKeyUseBatchSession, false, "użyj sesji wsadowej (batch). Domyślnie - klient używa sesji interaktywnej")
}

func GetUploaderConfig(vip *viper.Viper) UploaderConfig {
	return UploaderConfig{
		BatchSession: vip.GetBool(cfgKeyUseBatchSession),
	}
}
