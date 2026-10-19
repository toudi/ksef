package fs

import (
	"context"
	"ksef/core/certsdb"
)

func (fs *FSCertsStorage) UpsertSelfSignedCertificate(
	ctx context.Context,
	selfSignedCertData certsdb.UserCertificateData,
) error {
	selfSigned := fs.lookup(func(cert *fsCertificate) bool {
		return cert.SelfSigned
	})

	var addToIndex bool = selfSigned == nil
	var currentHash string
	var uid string
	var err error

	if selfSigned != nil {
		uid = selfSigned.UID
		currentHash = selfSigned.Hash()
	} else {
		uid, err = fs.generateUID()
		if err != nil {
			return err
		}
	}

	// this will make it CN to be repopulated by the next time certsDB is opened
	selfSigned = fromCertsDBCertificate(selfSignedCertData.Certificate)
	selfSigned.UID = uid

	if addToIndex {
		fs.addCertToIndex(selfSigned)
	} else {
		fs.replaceCertIndex(selfSigned, currentHash)
	}

	// overwrite private key bytes
	if err = fs.dataStorage.StorePrivateKeyBytes(ctx, uid, selfSignedCertData.PrivateKeyPEMBytes); err != nil {
		return err
	}
	// overwrite certificate DER bytes
	return fs.dataStorage.StoreCertificateBytes(ctx, uid, selfSignedCertData.DERBytes)
}
