package logging

import (
	"github.com/spf13/viper"
)

const (
	cfgKeyLogging = "logging"
)

var defaultLogging = map[string]string{
	"*": "info",
}

func LoggingConfig(vip *viper.Viper) map[string]string {
	loggingConfig := vip.GetStringMapString(cfgKeyLogging)
	if len(loggingConfig) > 0 {
		return loggingConfig
	}
	return defaultLogging
}
