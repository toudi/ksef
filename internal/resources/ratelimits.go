package resources

import (
	"context"
	"errors"
	"ksef/core/ratelimits"
	"path/filepath"

	ratelimitsfs "ksef/core/ratelimits/fs"

	"github.com/spf13/viper"
)

var (
	errNIPRequired          = errors.New("NIP is required to persist rate limits state")
	rateLimitsStateInstance map[string]ratelimits.StateStorage
)

func init() {
	rateLimitsStateInstance = make(map[string]ratelimits.StateStorage)
}

// RateLimitsState returns the (viper-free) storage used to persist the rate limiter
// state for the given NIP. Only two things used to come out of vip here - the
// invoices data root and the NIP - so the NIP is passed explicitly while the root
// comes from InvoiceStorageRoot (which deliberately does not resolve a NIP).
func RateLimitsState(vip *viper.Viper, NIP string) (ratelimits.StateStorage, error) {
	if NIP == "" {
		// filepath.Join(root, "") would silently point at <root>/<environment>
		return nil, errNIPRequired
	}

	mu.Lock()
	defer mu.Unlock()

	return rateLimitsStateLocked(vip, NIP), nil
}

func rateLimitsStateLocked(vip *viper.Viper, NIP string) ratelimits.StateStorage {
	if instance, exists := rateLimitsStateInstance[NIP]; exists {
		return instance
	}

	instance := ratelimitsfs.NewFSStorage(filepath.Join(InvoiceStorageRoot(vip), NIP))
	rateLimitsStateInstance[NIP] = instance

	OnClose(func(ctx context.Context, _ error) error {
		return instance.Close(context.WithoutCancel(ctx))
	})

	return instance
}
