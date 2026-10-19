package fs

import (
	"context"
	"errors"
	"ksef/core/certsdb"
	"ksef/internal/utils"
	"os"

	"github.com/samber/lo"
)

// this package implements the certificates database based on a filesystem directory
const (
	certificatesDir    = "certificates"
	certificatesDBFile = "certificates/certificates.yaml"
)

// FSCertsStorage is the metadata storage (certsdb.CertificateMetadataStorage):
// the in-memory index of the certificates.yaml file.
type FSCertsStorage struct {
	// environment id. the reason for it is that there's a single YAML file that
	// keeps all the certificates (yeah, that's a crappy solution..)
	// and so if we want to show certificates for some particular environment
	// we've got to filter them in-memory. hence this property here
	envId string

	// metadata storage keeps a reference to the data storage, because it is the
	// metadata that generates UID for the private key's and certificate.
	// in the filesystem implementation, they are actually shared (again, crappy idea, I know)
	// but in other implementations this does not need to be the case.
	dataStorage *FSDataStorage

	certs []*fsCertificate

	dirty     bool
	uidIndex  map[string]int
	hashIndex map[string]int
}

var _ certsdb.CertificateMetadataStorage = (*FSCertsStorage)(nil)

// NewFSCertStorage loads the certificates index from the YAML file and wires
// it with the data storage it should use for raw bytes.
func NewFSCertStorage(environmentId string, dataStorage *FSDataStorage) (*FSCertsStorage, error) {
	certsDB := &FSCertsStorage{
		envId:       environmentId,
		dataStorage: dataStorage,
		uidIndex:    make(map[string]int),
		hashIndex:   make(map[string]int),
	}

	dbFile, exists, err := utils.FileExists(certificatesDBFile)
	// isNotExist is fine
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	// if the error was that the file does not exist then there's no point
	// reading any certs as there won't be any
	if exists {
		defer dbFile.Close()
		var certificates []*fsCertificate
		if err = utils.ReadYAML(dbFile, &certificates); err != nil {
			return nil, errors.Join(ErrDecodingCertificates, err)
		}

		certsDB.certs = make([]*fsCertificate, 0, len(certificates))

		for index, cert := range certificates {
			if cert.postReadHook(environmentId) {
				certsDB.dirty = true
			}
			certsDB.certs = append(certsDB.certs, cert)
			certsDB.uidIndex[cert.UID] = index
			certsDB.hashIndex[cert.Hash()] = index
		}
	}

	return certsDB, nil
}

func (f *FSCertsStorage) Close(ctx context.Context) error {
	return f.save(ctx)
}

func (f *FSCertsStorage) addCertToIndex(cert *fsCertificate) {
	f.certs = append(f.certs, cert)
	f.uidIndex[cert.UID] = len(f.certs) - 1
	f.hashIndex[cert.CertificateHash.Hash()] = len(f.certs) - 1
	f.dirty = true
}

func (f *FSCertsStorage) replaceCertIndex(cert *fsCertificate, currentHash string) {
	index := f.uidIndex[cert.UID]
	f.certs[index] = cert
	delete(f.hashIndex, currentHash)
	f.hashIndex[cert.Hash()] = index
	f.dirty = true
}

func (f *FSCertsStorage) save(ctx context.Context) error {
	// TODO: move this functionality to a separate function, ideally with a confirm flag
	// for _, cert := range f.certs {
	// 	if cert.Expired() {
	// 		// logging.CertsDBLogger.With("id", cert.UID).Warn("certyfikat stracił ważność - usuwam plik")
	// 		os.Remove(cert.Filename())
	// 		cert.removed = true
	// 		f.dirty = true
	// 	}
	// }

	if !f.dirty {
		return nil
	}

	f.certs = lo.Filter(f.certs, func(item *fsCertificate, _ int) bool {
		return !item.removed
	})

	// just in case we altered NIP's assigned to certs
	for _, cert := range f.certs {
		cert.NIPRaw = cert.NIP
	}

	if err := utils.SaveYAML(f.certs, certificatesDBFile); err != nil {
		return err
	}
	f.dirty = false
	return nil
}
