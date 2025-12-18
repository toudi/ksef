package config

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	cfgKeyStatusWait        = "status.wait"
	cfgKeyStatusWaitTimeout = "status.wait.timeout"
	cfgKeyUPODownload       = "upo"
	cfgKeyUPOPdf            = "upo.pdf"
	cfgKeyUPOTimeout        = "upo.timeout"

	defaultWaitTimeout = time.Duration(5 * time.Minute)
)

type UPODownloaderConfig struct {
	Enabled      bool
	ConvertToPDF bool
	Timeout      time.Duration
}

type StatusCheckerConfig struct {
	Wait                bool
	WaitTimeout         time.Duration
	UPODownloaderConfig UPODownloaderConfig
}

func StatusCheckerFlags(flagSet *pflag.FlagSet) {
	flagSet.Bool(cfgKeyStatusWait, false, "czekaj na przetworzenie sesji")
	flagSet.Duration(cfgKeyStatusWaitTimeout, defaultWaitTimeout, "maksymalny czas oczekiwania na przetworzenie sesji")
	flagSet.Bool(cfgKeyUPODownload, false, "pobierz UPO po przetworzeniu sesji")
	flagSet.Bool(cfgKeyUPOPdf, false, "konwertuj UPO do PDF")
	flagSet.Duration(cfgKeyUPOTimeout, defaultWaitTimeout, "Maksymalny czas oczekiwania na rezultat pobrania UPO")
}

func GetStatusCheckerConfig(vip *viper.Viper) StatusCheckerConfig {
	var (
		waitTimeout        = defaultWaitTimeout
		upoDownloadTimeout = defaultWaitTimeout
	)

	if waitDur := vip.GetDuration(cfgKeyStatusWaitTimeout); waitDur > 0 {
		waitTimeout = waitDur
	}
	if upoDur := vip.GetDuration(cfgKeyUPOTimeout); upoDur > 0 {
		upoDownloadTimeout = upoDur
	}

	return StatusCheckerConfig{
		Wait:        vip.GetBool(cfgKeyStatusWait),
		WaitTimeout: waitTimeout,
		UPODownloaderConfig: UPODownloaderConfig{
			Enabled:      vip.GetBool(cfgKeyUPODownload),
			ConvertToPDF: vip.GetBool(cfgKeyUPOPdf),
			Timeout:      upoDownloadTimeout,
		},
	}
}
