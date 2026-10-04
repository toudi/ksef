package certsdb

import (
	"context"
	"iter"
)

func (c *certsDB) Certs(ctx context.Context) (iter.Seq[Certificate], error) {
	certs, err := c.storage.Certs(ctx)
	if err != nil {
		return nil, err
	}

	return func(yield func(Certificate) bool) {
		for _, cert := range certs {
			if !yield(cert) {
				return
			}
		}
	}, nil
}
