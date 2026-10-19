package fs

import (
	"context"
	"ksef/core/certsdb"
)

func (f *FSCertsStorage) lookup(matches func(cert *fsCertificate) bool) *fsCertificate {
	for _, cert := range f.certs {
		if matches(cert) {
			return cert
		}
	}

	return nil
}

func (f *FSCertsStorage) getCertByID(id string) (*fsCertificate, error) {
	certIdx, exists := f.uidIndex[id]
	if !exists {
		return nil, ErrUnknownCertificateUID
	}

	return f.certs[certIdx], nil
}

func (f *FSCertsStorage) containsHash(ctx context.Context, certHash certsdb.CertificateHash) (bool, error) {
	// convert public hash to inner hash
	fsCertHash := CertificateHash{
		EnvironmentId: string(certHash.Gateway),
		Usage:         certHash.Usage,
		ValidFrom:     certHash.ValidFrom,
		ValidTo:       certHash.ValidTo,
	}
	_, exists := f.hashIndex[fsCertHash.Hash()]
	return exists, nil
}

// GetCertificate performs a single lookup by certificate UID.
func (f *FSCertsStorage) GetCertificate(ctx context.Context, certId string) (*certsdb.Certificate, error) {
	cert, err := f.getCertByID(certId)
	if err != nil {
		return nil, err
	}

	result := cert.AsCertificate()
	return &result, nil
}
