package certsdb

import (
	"github.com/jaevor/go-nanoid"
)

var uidFactory func() string

func init() {
	uidFactory = nanoid.MustCustomASCII("abcdefgh0123456789", 5)
}

func (cdb *CertificatesDB) AddCert(handler func(newCert *Certificate) error) error {
	var collision = true
	var uid string
	for collision {
		uid = uidFactory()
		_, collision = cdb.uidIndex[uid]
	}
	newCert := &Certificate{
		UID: uid,
	}
	newCert.Environment = cdb.env

	if err := handler(newCert); err != nil {
		return err
	}

	cdb.certs = append(cdb.certs, newCert)
	cdb.index[newCert.Hash()] = len(cdb.certs) - 1
	cdb.uidIndex[uid] = len(cdb.certs) - 1
	cdb.dirty = true

	return nil
}

func (cdb *CertificatesDB) AddIfHashNotFound(hash CertificateHash, handler func(newCert *Certificate) error) error {
	_, exists := cdb.index[hash.Hash()]
	if exists {
		return nil
	}
	return cdb.AddCert(handler)
}

func (cdb *CertificatesDB) Upsert(
	lookup func(cert Certificate) bool,
	modify func(newCert *Certificate) error,
) error {
	var err error

	for _, cert := range cdb.certs {
		if lookup(*cert) {
			if err = modify(cert); err != nil {
				return err
			}
			return nil
		}
	}

	// cert is not found - let's create a new one
	return cdb.AddCert(modify)
}
