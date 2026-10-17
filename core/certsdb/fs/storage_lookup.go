package fs

import (
	"context"
	"ksef/core/certsdb"
	"os"
)

func (f *fsCertsStorage) lookup(matches func(cert *fsCertificate) bool) *fsCertificate {
	for _, cert := range f.certs {
		if matches(cert) {
			return cert
		}
	}

	return nil
}

func (f *fsCertsStorage) getCertByID(id string) (*fsCertificate, error) {
	certIdx, exists := f.uidIndex[id]
	if !exists {
		return nil, ErrUnknownCertificateUID
	}

	return f.certs[certIdx], nil
}

func (f *fsCertsStorage) GetCertificateBytes(ctx context.Context, certId string) ([]byte, error) {
	certFilename, err := f.certificateFilename(certId)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(certFilename)
}

func (f *fsCertsStorage) certificateFilename(certId string) (string, error) {
	cert, err := f.getCertByID(certId)
	if err != nil {
		return "", err
	}

	return cert.Filename(), nil
}

func (f *fsCertsStorage) containsHash(ctx context.Context, certHash certsdb.CertificateHash) (bool, error) {
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
