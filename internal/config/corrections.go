package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	cfgKeyCorrectionNumbering = "corrections.invoice_numbering"
)

type Corrections struct {
	Numbering string
}

func CorrectionsFlags(flags *pflag.FlagSet) error {
	flags.String(cfgKeyCorrectionNumbering, "FK/{count}/{year}", "Schemat numeracji faktur korygujących")

	return nil
}

func CorrectionsConfig(vip *viper.Viper) Corrections {
	return Corrections{
		Numbering: vip.GetString(cfgKeyCorrectionNumbering),
	}
}
