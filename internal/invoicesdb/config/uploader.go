package config

import (
	uploaderconfig "ksef/internal/invoicesdb/uploader/config"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	cfgKeyWaitForStatus = "upload.wait"
	cfgKeyWaitTimeout   = "upload.wait.timeout"
	cfgKeyDownloadUpo   = "upo"
	cfgKeyConvertUPOPdf = "upo.pdf"
)

func UploaderFlags(flagSet *pflag.FlagSet) {
	flagSet.Bool(cfgKeyWaitForStatus, false, "czekaj na zakończenie wysyłki")
	flagSet.Duration(cfgKeyWaitTimeout, time.Duration(0), "maksymalny czas oczekiwania na rezultat wysyłki")
	flagSet.Bool(cfgKeyDownloadUpo, false, "pobierz UPO po przetworzeniu sesji")
	flagSet.Bool(cfgKeyConvertUPOPdf, false, "konwertuj UPO do PDF")
}

func GetUploaderConfig(vip *viper.Viper) uploaderconfig.UploaderConfig {
	return uploaderconfig.UploaderConfig{
		WaitForStatus: vip.GetBool(cfgKeyWaitForStatus),
		WaitTimeout:   vip.GetDuration(cfgKeyWaitTimeout),
		DownloadUPO:   vip.GetBool(cfgKeyDownloadUpo),
		SaveUPOAsPDF:  vip.GetBool(cfgKeyConvertUPOPdf),
	}
}
