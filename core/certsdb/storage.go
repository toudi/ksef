package certsdb

import "context"

type CertificateStorage interface {
	Certs(ctx context.Context) ([]Certificate, error)
	Lookup(ctx context.Context, matches func(c *Certificate) bool) (*Certificate, error)

	// private key related functions
	GetPrivateKeyBytes(ctx context.Context, id string) ([]byte, error)

	// implementation may use this method whenever it wants to sync the internal state to persistance layer
	Close(ctx context.Context) error
}
