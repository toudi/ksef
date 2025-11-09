package types

import (
	"ksef/internal/client/v2/types/invoices"
	"time"
)

type QueryCriteria struct {
	DateFrom time.Time              `yaml:"dateFrom"`
	DateTo   time.Time              `yaml:"dateTo,omitempty"`
	DateType invoices.DateRangeType `yaml:"dateType,omitempty"`
}

type SubjectSyncConfig struct {
	SubjectType        invoices.SubjectType `yaml:"type"`
	Incremental        bool                 `yaml:"incremental,omitempty"`
	LastKnownTimestamp time.Time            `yaml:"last-known-timestamp,omitempty"`
}

type SyncConfig struct {
	QueryCriteria QueryCriteria       `yaml:"query,omitempty"`
	Subjects      []SubjectSyncConfig `yaml:"subject-types,omitempty"`
}
