package fs

import (
	"bytes"
	"context"
	"ksef/core/certsdb"
	"ksef/internal/utils"
)

func (fs *fsCertsStorage) StoreCSRRequest(
	ctx context.Context,
	cert certsdb.Certificate,
) (string, error) {
	// in this legacy implementation the private key and certificate share the same
	// uid. looking back, it's a crappy implementation :/

	fs.certs = append(fs.certs, &fsCertificate{
		CertificateHash: CertificateHash{
			Usage:         cert.Usage,
			EnvironmentId: string(cert.Gateway),
		},
		UID:     cert.PrivateKeyId,
		NIP:     cert.NIP,
		CSRData: cert.CSRData,
	})
	fs.uidIndex[cert.PrivateKeyId] = len(fs.certs) - 1
	fs.dirty = true
	return cert.PrivateKeyId, nil
}

func (fs *fsCertsStorage) StoreMinistryCertificate(
	ctx context.Context,
	certData certsdb.MinistryCertificateData,
) (bool, error) {
	// first, look up if there already is a certificate with the given hash
	if exists, _ := fs.containsHash(ctx, certData.Certificate.CertificateHash); exists {
		return false, nil
	}

	fsCert := fromCertsDBCertificate(certData.Certificate)

	uid, err := fs.generateUID()
	if err != nil {
		return false, err
	}

	fsCert.UID = uid

	if err := utils.SaveBufferToFile(*bytes.NewBuffer(certData.DERBytes), fsCert.Filename()); err != nil {
		return false, err
	}

	fs.addCertToIndex(fsCert)

	return true, nil
}

func (fs *fsCertsStorage) StoreUserCertificate(
	ctx context.Context,
	certData certsdb.UserCertificateData,
) (bool, error) {
	// first, look up if there already is a certificate with the given hash
	if exists, _ := fs.containsHash(ctx, certData.Certificate.CertificateHash); exists {
		return false, nil
	}

	fsCert := fromCertsDBCertificate(certData.Certificate)

	uid, err := fs.generateUID()
	if err != nil {
		return false, err
	}

	fsCert.UID = uid

	// store primary key
	if err = fs.storePrivateKeyBytes(ctx, fsCert, certData.PrivateKeyPEMBytes); err != nil {
		return false, err
	}
	// store certificate data
	if err = fs.storeCertificateBytes(ctx, fsCert, certData.DERBytes); err != nil {
		return false, err
	}
	fs.addCertToIndex(fsCert)
	return true, nil
}

func (fs *fsCertsStorage) storePrivateKeyBytes(ctx context.Context, cert *fsCertificate, keyBytes []byte) error {
	privateKeyFilename := fs.privateKeyFilename(cert)
	if err := fs.encryptPrivateKeyBytesToFile(cert.UID, keyBytes, privateKeyFilename); err != nil {
		return err
	}
	return nil
}

func (fs *fsCertsStorage) storeCertificateBytes(ctx context.Context, cert *fsCertificate, certBytes []byte) error {
	certFilename := cert.Filename()
	if err := utils.SaveBufferToFile(*bytes.NewBuffer(certBytes), certFilename); err != nil {
		return err
	}
	return nil
}

func (fs *fsCertsStorage) addCertToIndex(cert *fsCertificate) {
	fs.certs = append(fs.certs, cert)
	fs.uidIndex[cert.UID] = len(fs.certs) - 1
	fs.hashIndex[cert.CertificateHash.Hash()] = len(fs.certs) - 1
	fs.dirty = true
}

func (fs *fsCertsStorage) replaceCertIndex(cert *fsCertificate, currentHash string) {
	index := fs.uidIndex[cert.UID]
	fs.certs[index] = cert
	delete(fs.hashIndex, currentHash)
	fs.hashIndex[cert.Hash()] = index
	fs.dirty = true
}
