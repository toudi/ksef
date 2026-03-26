package certsdb

import (
	"iter"
	"slices"
)

func (cdb *CertificatesDB) ForEach(
	lookup func(cert Certificate) bool,
	update func(cert *Certificate) error,
) error {
	found := 0
	var err error

	for _, cert := range cdb.certs {
		if lookup(*cert) {
			if err = update(cert); err != nil {
				return err
			}
			found += 1
		}
	}

	if found > 0 {
		cdb.dirty = true
	}

	return nil
}

func (cdb *CertificatesDB) Filter(
	lookup func(cert Certificate) bool,
) iter.Seq[Certificate] {
	return func(yield func(Certificate) bool) {
		for _, cert := range cdb.certs {
			if lookup(*cert) {
				if !yield(*cert) {
					return
				}
			}
		}
	}
}

func (cdb *CertificatesDB) FetchUIDsByNIP(nip string) (uids []string) {
	for _, cert := range cdb.certs {
		if slices.Contains(cert.NIP, nip) {
			uids = append(uids, cert.UID)
		}
	}

	return uids
}
