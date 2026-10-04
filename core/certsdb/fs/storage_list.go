package fs

import (
	"context"
	"ksef/core/certsdb"

	"github.com/samber/lo"
)

func (f *fsCertsStorage) Certs(ctx context.Context) ([]certsdb.Certificate, error) {
	return lo.Map(f.certs, func(c *fsCertificate, _ int) certsdb.Certificate {
		return c.AsCertificate()
	}), nil
}
