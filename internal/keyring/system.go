package keyring

import (
	"ksef/internal/logging"
	"strings"

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
	logging.KeyringLogger.Debug("delete key", "service", serviceName(bucket, key), "user", nip)
	return zalandoKeyring.Delete(serviceName(bucket, key), nip)
}

func (s *SystemKeyring) Get(bucket string, nip string, key string) (string, error) {
	logging.KeyringLogger.Debug("get key value", "service", serviceName(bucket, key), "user", nip)

	value, err := zalandoKeyring.Get(serviceName(bucket, key), nip)
	if err == zalandoKeyring.ErrNotFound {
		err = ErrNotFound
	} else if err != nil && strings.HasPrefix(err.Error(), "failed to unlock correct collection") {
		err = ErrPermissionDenied
	}

	return value, err
}

func (s *SystemKeyring) Set(bucket string, nip string, key string, contents string) error {
	logging.KeyringLogger.Debug("set key", "service", serviceName(bucket, key), "user", nip)

	return zalandoKeyring.Set(serviceName(bucket, key), nip, contents)
}

func (s *SystemKeyring) Close() error {
	return nil
}
