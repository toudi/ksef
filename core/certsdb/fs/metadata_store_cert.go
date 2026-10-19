package fs

import (
	"context"
	"ksef/core/certsdb"
)

func (fs *FSCertsStorage) StoreCSRRequest(
	ctx context.Context,
	csrData certsdb.CSRRequestData,
) (string, error) {
	// in the filesystem implementation the private key and the certificate
	// share the same uid
	uid, err := fs.generateUID()
	if err != nil {
		return "", err
	}
	csrData.Certificate.PrivateKeyId = uid

	fsCert := fromCertsDBCertificate(csrData.Certificate)
	fsCert.UID = uid
	fs.addCertToIndex(fsCert)

	if err = fs.dataStorage.StorePrivateKeyBytes(ctx, uid, csrData.PrivateKeyPEMBytes); err != nil {
		return "", err
	}

	return uid, nil
}

func (fs *FSCertsStorage) StoreMinistryCertificate(
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
	fs.addCertToIndex(fsCert)

	// store certificate data
	return true, fs.dataStorage.StoreCertificateBytes(ctx, uid, certData.DERBytes)
}

func (fs *FSCertsStorage) StoreUserCertificate(
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
	fs.addCertToIndex(fsCert)

	// store private key
	if err = fs.dataStorage.StorePrivateKeyBytes(ctx, uid, certData.PrivateKeyPEMBytes); err != nil {
		return false, err
	}
	// store certificate data
	return true, fs.dataStorage.StoreCertificateBytes(ctx, uid, certData.DERBytes)
}
