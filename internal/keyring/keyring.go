package keyring

import (
	"errors"
	"ksef/internal/config"

	"github.com/spf13/viper"
)

type Keyring interface {
	Get(bucket string, nip string, key string) (string, error)
	Set(bucket string, nip string, key string, contents string) error
	Delete(bucket string, nip string, key string) error
	Close() error
}

var ErrNotFound = errors.New("key not found")

func NewKeyring(vip *viper.Viper) (Keyring, error) {
	cfg, err := config.KeyringConfig(vip)
	if err != nil {
		return nil, err
	}

	if cfg.File != nil {
		return NewFileBasedKeyring(cfg.File)
	}

	return NewSystemKeyring(), nil
}
