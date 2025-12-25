package config

import (
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	cfgKeyUseBatchSession     = "upload.batch"
	cfgKeyBatchSessionWorkdir = "upload.batch.workdir"
)

type UploaderConfig struct {
	BatchSession bool
	BatchWorkdir string
}

func UploaderFlags(flagSet *pflag.FlagSet) {
	flagSet.Bool(cfgKeyUseBatchSession, false, "użyj sesji wsadowej (batch). Domyślnie - klient używa sesji interaktywnej")
	flagSet.String(cfgKeyBatchSessionWorkdir, os.TempDir(), "katalog roboczy sesji wsadowej")
}

func GetUploaderConfig(vip *viper.Viper) UploaderConfig {
	batchWorkdir := os.TempDir()
	if workDir := vip.GetString(cfgKeyBatchSessionWorkdir); workDir != "" {
		batchWorkdir = workDir
	}
	return UploaderConfig{
		BatchSession: vip.GetBool(cfgKeyUseBatchSession),
		BatchWorkdir: batchWorkdir,
	}
}
