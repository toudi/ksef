package resources

import (
	"context"
	kr "ksef/core/keyring"

	"github.com/spf13/viper"
)

var keyringInstance kr.Keyring

// Keyring returns the process-wide keyring, creating it on first use.
// The file-based engine may prompt the user for the keyring password -
// at most once per run, regardless of how many components request it.
func Keyring(vip *viper.Viper) (kr.Keyring, error) {
	if keyringInstance != nil {
		return keyringInstance, nil
	}

	ring, err := kr.NewKeyring(vip)
	if err != nil {
		return nil, err
	}
	keyringInstance = ring
	destructors = append(destructors, func(ctx context.Context) error {
		return ring.Close()
	})

	return ring, nil
}
