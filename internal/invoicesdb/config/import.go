package config

import (
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

type ImportConfig struct {
	AutoCorrection autoCorrectionConfig
	AutoUpload     bool
}

func ImportFlags(flagSet *pflag.FlagSet) {
	flagSet.Bool(cfgKeyAutoCorrection, false, "automatycznie wystawiaj korekty faktur")
	flagSet.String(cfgKeyCorrectionNumbering, "FK/{count}/{year}", "Schemat numeracji faktur korygujących")
	flagSet.BoolP(cfgKeyAutoUpload, "u", false, "automatycznie wyślij faktury po zakończonym imporcie")
}

func GetImportConfig(vip *viper.Viper) ImportConfig {
	return ImportConfig{
		AutoCorrection: autoCorrectionConfig{
			Enabled:         vip.GetBool(cfgKeyAutoCorrection),
			NumberingScheme: vip.GetString(cfgKeyCorrectionNumbering),
		},
		AutoUpload: vip.GetBool(cfgKeyAutoUpload),
	}
}
