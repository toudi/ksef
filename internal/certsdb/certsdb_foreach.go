package certsdb

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

func (cdb *CertificatesDB) FetchUIDsByNIP(nip string) (uids []string) {
	for _, cert := range cdb.certs {
		if cert.NIP == nip {
			uids = append(uids, cert.UID)
		}
	}

	return uids
}
