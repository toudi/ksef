package certsdb

import "errors"

var errUnableToFindCertificate = errors.New("unable to find certificate")

func (c *CertificatesDB) GetByUID(uid string) (Certificate, error) {
	certIndex, exists := c.index[uid]
	if !exists {
		return Certificate{}, ErrCertificateNotFound
	}

	return *(c.certs[certIndex]), nil
}

func (c *CertificatesDB) Lookup(matches func(c *Certificate) bool) (*Certificate, error) {
	for _, certificate := range c.certs {
		if matches(certificate) {
			return certificate, nil
		}
	}

	return nil, errUnableToFindCertificate
}
