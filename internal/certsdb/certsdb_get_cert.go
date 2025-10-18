package certsdb

func (c *CertificatesDB) GetByUID(uid string) (Certificate, error) {
	certIndex, exists := c.index[uid]
	if !exists {
		return Certificate{}, ErrCertificateNotFound
	}

	return *(c.certs[certIndex]), nil
}
