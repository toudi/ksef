package ksef

import "time"

// HTTPConfig carries the HTTP transport tuning needed by the KSeF client and
// its downloaders.
type HTTPConfig struct {
	// MaxRetries is the number of attempts a request is retried.
	MaxRetries int
	// PollWaitTime is the delay used when polling long-running operations.
	PollWaitTime time.Duration
}
