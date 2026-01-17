package certsdb

import (
	"errors"
	"ksef/internal/runtime"
	"ksef/internal/utils"
	"os"
	"path"

	"github.com/samber/lo"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	ErrCertificateNotFound  = errors.New("unable to find a certificate suitable for the selected usage")
	errDecodingCertificates = errors.New("unable to decode certificates array from file")
)

type (
	Usage string
)

func (u Usage) Description() string {
	if u == UsageAuthentication {
		return "certyfikat autoryzacyjny"
	} else if u == UsageOffline {
		return "certyfikat offline"
	}
	panic("nieoczekiwana wartość Usage")
}

const (
	certificatesDBFile = "certificates/certificates.yaml"
	// certyfikaty RSA używane przez ministerstwo finansów
	UsageTokenEncryption        Usage = "KsefTokenEncryption"
	UsageSymmetricKeyEncryption Usage = "SymmetricKeyEncryption"
)

type CertificatesDB struct {
	certs []*Certificate
	// whether the db requires to be saved
	dirty bool
	// map between certificate hash and it's position in the array
	index map[string]int
	// map between certificate UID and it's position in the array - required only for adding new certificates
	uidIndex map[string]int
	// used during opening so that all of the certs inserted inherit this value
	env runtime.Gateway
	// pointer to viper so that we can read preferred cert ID
	vip *viper.Viper
}

func (cdb *CertificatesDB) Certs() []*Certificate {
	return cdb.certs
}

func (cdb *CertificatesDB) Save() error {
	// TODO: move this functionality to a separate function, ideally with a confirm flag
	// for _, cert := range cdb.certs {
	// 	if cert.Expired() {
	// 		// logging.CertsDBLogger.With("id", cert.UID).Warn("certyfikat stracił ważność - usuwam plik")
	// 		os.Remove(cert.Filename())
	// 		cert.removed = true
	// 		cdb.dirty = true
	// 	}
	// }

	if !cdb.dirty {
		return nil
	}

	cdb.certs = lo.Filter(cdb.certs, func(item *Certificate, _ int) bool {
		return !item.removed
	})

	// just in case we altered NIP's assigned to certs
	for _, cert := range cdb.certs {
		cert.NIPRaw = cert.NIP
	}

	targetFile, err := os.Create(certificatesDBFile)
	if err != nil {
		return err
	}
	defer func() {
		targetFile.Close()
		cdb.dirty = false
	}()
	encoder := yaml.NewEncoder(targetFile)
	return encoder.Encode(cdb.certs)
}

func OpenOrCreate(vip *viper.Viper) (*CertificatesDB, error) {
	environment := runtime.GetGateway(vip)

	certificatesDB := CertificatesDB{
		index:    make(map[string]int),
		uidIndex: make(map[string]int),
		env:      environment,
		vip:      vip,
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
		var certificates []*Certificate
		if err = utils.ReadYAML(dbFile, &certificates); err != nil {
			return nil, errors.Join(errDecodingCertificates, err)
		}

		certificatesDB.certs = certificates

		for index, cert := range certificates {
			if cert.postReadHook(environment) {
				certificatesDB.dirty = true
			}
			certificatesDB.index[cert.Hash()] = index
			certificatesDB.uidIndex[cert.UID] = index
		}
	}

	return &certificatesDB, nil
}
