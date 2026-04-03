package keyring

import (
	"errors"
)

type Keyring interface {
	Get(bucket string, nip string, key string) (string, error)
	Set(bucket string, nip string, key string, contents string) error
	Delete(bucket string, nip string, key string) error
	Close() error
}

var ErrNotFound = errors.New("key not found")
var ErrPermissionDenied = errors.New("permission denied to read key")
