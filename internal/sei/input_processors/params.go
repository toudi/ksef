package inputprocessors

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type csvConfig struct {
	Delimiter              string
	EncodingConversionFile string
}

type xlsxConfig struct {
	SheetName string
}

type InputProcessorConfig struct {
	CSV         csvConfig
	XLSX        xlsxConfig
	Generator   string
	OfflineMode bool
}

const (
	cfgKeyCSVDelimiter = "csv.delimiter"
	cfgKeyCSVEncoding  = "csv.encoding"
	cfgKeySheetName    = "xlsx.sheet"
	cfgKeyGenerator    = "generator"
	cfgKeyOffline      = "offline"
)

func GeneratorFlags(flags *pflag.FlagSet) {
	flags.StringP(cfgKeyCSVDelimiter, "d", ",", "łańcuch znaków rozdzielający pola (tylko dla CSV)")
	flags.StringP(cfgKeyCSVEncoding, "e", "", "użyj pliku z konwersją znaków (tylko dla CSV)")
	flags.StringP(cfgKeySheetName, "s", "", "Nazwa skoroszytu do przetworzenia (tylko dla XLSX)")
	flags.StringP(cfgKeyGenerator, "g", "fa-3_1.0", "nazwa generatora")
	flags.Bool(cfgKeyOffline, false, "oznacz faktury jako generowane w trybie offline")
}

func GetInputProcessorConfig(vip *viper.Viper) InputProcessorConfig {
	return InputProcessorConfig{
		CSV: csvConfig{
			Delimiter:              vip.GetString(cfgKeyCSVDelimiter),
			EncodingConversionFile: vip.GetString(cfgKeyCSVEncoding),
		},
		XLSX: xlsxConfig{
			SheetName: viper.GetString(cfgKeySheetName),
		},
		Generator:   viper.GetString(cfgKeyGenerator),
		OfflineMode: viper.GetBool(cfgKeyOffline),
	}
}
