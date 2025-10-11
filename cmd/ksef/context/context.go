package context

import (
	"errors"
)

type ContextKey string

const (
	KeyEnvironment ContextKey = "environment"
)

var (
	errInvalidContext = errors.New("unable to retrieve app context")
)
