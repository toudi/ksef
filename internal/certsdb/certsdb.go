package certsdb

import (
	"errors"
	"ksef/internal/environment"
	"os"
	"path"
	"slices"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

var ErrCertificateNotFound = errors.New("unable to find a certificate suitable for the selected usage")

type Usage string

const (
	certificatesDBFile                = "certificates/certificates.yaml"
	UsageTokenEncryption        Usage = "KsefTokenEncryption"
	UsageSymmetricKeyEncryption Usage = "SymmetricKeyEncryption"
)

type CertificateFile struct {
	Environment environment.Environment `yaml:"environment"`
	Usage       []Usage                 `yaml:"usage"`
	DERFile     string                  `yaml:"der-file"`
	PEMFile     string                  `yaml:"pem-file"`
}

type CertificatesDB struct {
	certs []CertificateFile
	// whether the db requires to be saved
	dirty bool
}

func (cdb *CertificatesDB) GetByUsage(usage Usage) (CertificateFile, error) {
	for _, cert := range cdb.certs {
		if slices.Contains(cert.Usage, usage) {
			return cert, nil
		}
	}

	return CertificateFile{}, ErrCertificateNotFound
}

func (cdb *CertificatesDB) Save() error {
	if !cdb.dirty {
		return nil
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

func OpenOrCreate(environment environment.Environment) (*CertificatesDB, error) {
	var certificatesDB CertificatesDB

	if err := os.MkdirAll(path.Dir(certificatesDBFile), 0775); err != nil {
		return nil, err
	}
	dbFile, err := os.Open(certificatesDBFile)
	// isNotExist is fine
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	// if the error was that the file does not exist then there's no point
	// reading any certs as there won't be any
	if err == nil {
		defer dbFile.Close()
		var certificates []CertificateFile
		decoder := yaml.NewDecoder(dbFile)
		if err := decoder.Decode(&certificates); err != nil {
			return nil, err
		}

		// load up only the certificates that belong to the selected environment
		certificatesDB.certs = lo.Filter(certificates, func(c CertificateFile, _ int) bool {
			return c.Environment == environment
		})
	}

	return &certificatesDB, nil
}
