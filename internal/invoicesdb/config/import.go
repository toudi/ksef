package config

import (
	uploaderconfig "ksef/internal/invoicesdb/uploader/config"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	cfgKeyAutoCorrection      = "auto-correction"
	cfgKeyCorrectionNumbering = "corrections.numbering"
	cfgKeyAutoUpload          = "upload"
)

type autoCorrectionConfig struct {
	Enabled         bool
	NumberingScheme string
}

type UploadConfig struct {
	Enabled        bool
	UploaderConfig uploaderconfig.UploaderConfig
}

type ImportConfig struct {
	AutoCorrection autoCorrectionConfig
	Upload         UploadConfig
}

func ImportFlags(flagSet *pflag.FlagSet) {
	flagSet.Bool(cfgKeyAutoCorrection, false, "automatycznie wystawiaj korekty faktur")
	flagSet.String(cfgKeyCorrectionNumbering, "FK/{count}/{year}", "Schemat numeracji faktur korygujących")
	flagSet.BoolP(cfgKeyAutoUpload, "u", false, "automatycznie wyślij faktury po zakończonym imporcie")
	UploaderFlags(flagSet)
}

func GetImportConfig(vip *viper.Viper) ImportConfig {
	return ImportConfig{
		AutoCorrection: autoCorrectionConfig{
			Enabled:         vip.GetBool(cfgKeyAutoCorrection),
			NumberingScheme: vip.GetString(cfgKeyCorrectionNumbering),
		},
		Upload: UploadConfig{
			Enabled:        vip.GetBool(cfgKeyAutoUpload),
			UploaderConfig: GetUploaderConfig(vip),
		},
	}
}
