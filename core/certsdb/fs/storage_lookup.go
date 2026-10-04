package fs

import (
	"context"
	"ksef/core/certsdb"
)

func (f *fsCertsStorage) Lookup(ctx context.Context, matches func(c *certsdb.Certificate) bool) (*certsdb.Certificate, error) {
	return nil, nil
}

func (f *fsCertsStorage) getCertByID(id string) (*fsCertificate, error) {
	certIdx, exists := f.uidIndex[id]
	if !exists {
		return nil, ErrUnknownCertificateUID
	}

	return f.certs[certIdx], nil
}
