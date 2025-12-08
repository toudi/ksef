package config

import (
	"time"
)

type UploaderConfig struct {
	WaitForStatus bool
	WaitTimeout   time.Duration
	DownloadUPO   bool
	SaveUPOAsPDF  bool
}
