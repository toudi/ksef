package interfaces

import "time"

type TokenRetrieverFunc func(timeout ...time.Duration) (string, error)
