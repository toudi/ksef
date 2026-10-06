package fs

import (
	"context"
	"errors"
	"ksef/core/certsdb"
	"ksef/core/keyring"
	"ksef/internal/runtime"
	"ksef/internal/utils"
	"os"
	"path"

	"github.com/samber/lo"
	"github.com/spf13/viper"
)

// this package implements CertificatesDB based on a filesystem directory
const certificatesDBFile = "certificates/certificates.yaml"

type fsCertsStorage struct {
	certs []*fsCertificate
	kr    keyring.Keyring

	dirty    bool
	index    map[string]int
	uidIndex map[string]int
	env      string
	vip      *viper.Viper
}

var _ certsdb.CertificateStorage = (*fsCertsStorage)(nil)

func Init(vip *viper.Viper) (*fsCertsStorage, error) {
	environment := runtime.GetEnvironmentId(vip)
	// initialize keyring instance
	kr, err := keyring.NewKeyring(vip)
	if err != nil {
		return nil, err
	}

	certificatesDB := fsCertsStorage{
		uidIndex: make(map[string]int),
		env:      environment,
		vip:      vip,
		kr:       kr,
	}

	if err := os.MkdirAll(path.Dir(certificatesDBFile), 0775); err != nil {
		return nil, err
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

		certificatesDB.certs = make([]*fsCertificate, 0, len(certificates))

		for index, cert := range certificates {
			if cert.postReadHook(environment) {
				certificatesDB.dirty = true
			}
			certificatesDB.certs = append(certificatesDB.certs, cert)
			certificatesDB.uidIndex[cert.UID] = index
		}
	}

	return &certificatesDB, nil
}

func (f *fsCertsStorage) Close(ctx context.Context) error {
	// this is required, since if somebody would be using file-based keyring we
	// actually have to save it
	if f.kr != nil {
		if err := f.kr.Close(); err != nil {
			return err
		}
	}
	return f.save(ctx)
}

func (f *fsCertsStorage) save(ctx context.Context) error {
	// TODO: move this functionality to a separate function, ideally with a confirm flag
	// for _, cert := range cdb.certs {
	// 	if cert.Expired() {
	// 		// logging.CertsDBLogger.With("id", cert.UID).Warn("certyfikat stracił ważność - usuwam plik")
	// 		os.Remove(cert.Filename())
	// 		cert.removed = true
	// 		cdb.dirty = true
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
