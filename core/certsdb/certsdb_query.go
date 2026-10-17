package certsdb

import (
	"context"
	"slices"

	"github.com/samber/lo"
)

func (c *CertsDB) FilterCertificates(ctx context.Context, checker func(c Certificate) bool) ([]Certificate, error) {
	certificates, err := c.storage.Certs(ctx)
	if err != nil {
		return nil, err
	}
	return lo.Filter(certificates, func(c Certificate, _ int) bool {
		return checker(c)
	}), nil
}

func (c *CertsDB) FetchUIDsByNIP(ctx context.Context, nip string) (uids []string, err error) {
	certs, err := c.storage.Certs(ctx)
	if err != nil {
		return nil, err
	}

	for _, cert := range certs {
		if slices.Contains(cert.NIP, nip) {
			uids = append(uids, cert.UID)
		}
	}

	return uids, nil
}
