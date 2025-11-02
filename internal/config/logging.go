package config

import (
	"github.com/spf13/viper"
)

const (
	cfgKeyLogging = "logging"
)

func LoggingConfig(vip *viper.Viper) map[string]string {
	return vip.GetStringMapString(cfgKeyLogging)
}
