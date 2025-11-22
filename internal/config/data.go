package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// all data files, i.e. downloaded invoices, certs etc are going
	// to live underneath data dir
	cfgDataDir = "data-dir"
)

func DataDir(vip *viper.Viper) string {
	return vip.GetString(cfgDataDir)
}

func DataDirFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(cfgDataDir, "data", "katalog danych")
}
