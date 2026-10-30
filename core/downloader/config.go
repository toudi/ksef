package downloader

import (
	"ksef/core/ksef"
	"time"
)

type Config struct {
	Incremental   bool
	PDF           bool
	UseExportMode bool
	UseSmartMode  bool
	SubjectTypes  []ksef.SubjectType
	DateType      ksef.DateRangeType
	StartDate     time.Time
	EndDate       *time.Time
	PageSize      int
}
