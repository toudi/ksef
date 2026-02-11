package config

import (
	uploaderconfig "ksef/internal/invoicesdb/uploader/config"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	cfgKeyAutoUpload = "upload"
	cfgKeyOffline    = "offline"
)

type UploadConfig struct {
	Enabled        bool
	UploaderConfig uploaderconfig.UploaderConfig
}

type ImportConfig struct {
	Upload  UploadConfig
	Offline bool
}

func ImportFlags(flagSet *pflag.FlagSet) {
	flagSet.BoolP(cfgKeyAutoUpload, "u", false, "automatycznie wyślij faktury po zakończonym imporcie")
	flagSet.Bool(cfgKeyOffline, false, "oznacz faktury jako generowane w trybie offline")
}

func GetImportConfig(vip *viper.Viper) ImportConfig {
	return ImportConfig{
		Upload: UploadConfig{
			Enabled:        vip.GetBool(cfgKeyAutoUpload),
			UploaderConfig: uploaderconfig.GetUploaderConfig(vip),
		},
		Offline: vip.GetBool(cfgKeyOffline),
	}
}
