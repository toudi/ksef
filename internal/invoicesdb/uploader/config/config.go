package config

import (
	"time"
)

type UPODownloaderConfig struct {
	Enabled      bool
	ConvertToPDF bool
	Timeout      time.Duration
}
type UploaderConfig struct {
	WaitForStatus bool
	WaitTimeout   time.Duration
	BatchSession  bool
	UPODownloader UPODownloaderConfig
}
