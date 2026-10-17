package fs

import (
	"context"
	"slices"
)

func (fs *fsCertsStorage) AddNIPToCertificate(
	ctx context.Context,
	certId string,
	NIP string,
) error {
	cert, err := fs.getCertByID(certId)
	if err != nil {
		return err
	}
	if !slices.Contains(cert.NIP, NIP) {
		cert.NIP = append(cert.NIP, NIP)
		fs.dirty = true
	}
	return nil
}

func (fs *fsCertsStorage) RemoveNIPFromCertificate(
	ctx context.Context,
	certId string,
	NIP string,
) error {
	cert, err := fs.getCertByID(certId)
	if err != nil {
		return err
	}
	if slices.Contains(cert.NIP, NIP) {
		cert.NIP = slices.DeleteFunc(cert.NIP, func(nip string) bool {
			return nip == NIP
		})
		fs.dirty = true
	}
	return nil
}
