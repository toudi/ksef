package config

import (
	uploaderconfig "ksef/internal/invoicesdb/uploader/config"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	cfgKeyUseBatchSession = "upload.batch"
	cfgKeyWaitForStatus   = "upload.wait"
	cfgKeyWaitTimeout     = "upload.wait.timeout"
	cfgKeyUPODownload     = "upo"
	cfgKeyUPOPdf          = "upo.pdf"
	cfgKeyUPOTimeout      = "upo.timeout"
	defaultWaitTimeout    = time.Duration(5 * time.Minute)
)

func UploaderFlags(flagSet *pflag.FlagSet) {
	flagSet.Bool(cfgKeyUseBatchSession, false, "użyj sesji wsadowej (batch). Domyślnie - klient używa sesji interaktywnej")
	flagSet.Bool(cfgKeyWaitForStatus, false, "czekaj na zakończenie wysyłki")
	flagSet.Duration(cfgKeyWaitTimeout, defaultWaitTimeout, "maksymalny czas oczekiwania na rezultat wysyłki")
	flagSet.Bool(cfgKeyUPODownload, false, "pobierz UPO po przetworzeniu sesji")
	flagSet.Bool(cfgKeyUPOPdf, false, "konwertuj UPO do PDF")
	flagSet.Duration(cfgKeyUPOTimeout, defaultWaitTimeout, "Maksymalny czas oczekiwania na rezultat pobrania UPO")
}

func GetUploaderConfig(vip *viper.Viper) uploaderconfig.UploaderConfig {
	var (
		statusWaitTimeout  = defaultWaitTimeout
		upoDownloadTimeout = defaultWaitTimeout
	)
	if waitDur := vip.GetDuration(cfgKeyWaitTimeout); waitDur > 0 {
		statusWaitTimeout = waitDur
	}
	if upoDur := vip.GetDuration(cfgKeyUPOTimeout); upoDur > 0 {
		upoDownloadTimeout = upoDur
	}

	return uploaderconfig.UploaderConfig{
		WaitForStatus: vip.GetBool(cfgKeyWaitForStatus),
		WaitTimeout:   statusWaitTimeout,
		BatchSession:  vip.GetBool(cfgKeyUseBatchSession),
		UPODownloader: uploaderconfig.UPODownloaderConfig{
			Enabled:      vip.GetBool(cfgKeyUPODownload),
			ConvertToPDF: vip.GetBool(cfgKeyUPOPdf),
			Timeout:      upoDownloadTimeout,
		},
	}
}
