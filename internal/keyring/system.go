package keyring

import (
	"ksef/internal/logging"

	zalandoKeyring "github.com/zalando/go-keyring"
)

type SystemKeyring struct{}

func serviceName(bucket string, key string) string {
	return AppPrefix + "-" + bucket + "-" + key
}

func NewSystemKeyring() *SystemKeyring {
	logging.KeyringLogger.Debug("system keyring initialization")
	return &SystemKeyring{}
}

func (s *SystemKeyring) Delete(bucket string, nip string, key string) error {
	return zalandoKeyring.Delete(serviceName(bucket, key), nip)
}

func (s *SystemKeyring) Get(bucket string, nip string, key string) (string, error) {
	value, err := zalandoKeyring.Get(serviceName(bucket, key), nip)
	if err == zalandoKeyring.ErrNotFound {
		err = ErrNotFound
	}
	return value, err
}

func (s *SystemKeyring) Set(bucket string, nip string, key string, contents string) error {
	return zalandoKeyring.Set(serviceName(bucket, key), nip, contents)
}

func (s *SystemKeyring) Close() error {
	return nil
}
