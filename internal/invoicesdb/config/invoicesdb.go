package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	cfgKeyDataDir = "data-dir"
)

type InvoicesDBConfig struct {
	Root string
}

func InvoicesDBFlags(flagSet *pflag.FlagSet) {
	flagSet.String(cfgKeyDataDir, "data", "katalog bazy faktur")
}

func GetInvoicesDBConfig(vip *viper.Viper) InvoicesDBConfig {
	return InvoicesDBConfig{
		Root: viper.GetString(cfgKeyDataDir),
	}
}

func SetDataDir(vip *viper.Viper, dataDir string) {
	vip.Set(cfgKeyDataDir, dataDir)
}
