package certsdb

import (
	"encoding/base64"
	"fmt"
	"io"
	"ksef/internal/environment"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/samber/lo"
)

func (cdb *CertificatesDB) AddCertificate(base64Der string, environment environment.Environment, usage []Usage) error {
	// example:
	// ksef-test.mf.gov.pl-KSeFTokenEncryption
	certBaseFilename := fmt.Sprintf(
		"%s-%s",
		environment,
		strings.Join(
			lo.Map(usage, func(elem Usage, _ int) string { return string(elem) }), "-",
		),
	)

	// sort the usage slice so that we don't end up writing the same certificate twice
	// just because it's usage array will be reversed
	slices.SortFunc(usage, func(e1, e2 Usage) int {
		return strings.Compare(string(e1), string(e2))
	})

	PEMFilename := path.Join(path.Dir(certificatesDBFile), certBaseFilename+".pem")
	DERFilename := path.Join(path.Dir(certificatesDBFile), certBaseFilename+".der")
	// so the base64-encoded content is actually the essense of PEM so let's use a nifty hack to save it
	pemFile, err := os.Create(PEMFilename)
	if err != nil {
		return err
	}
	defer pemFile.Close()
	if _, err = fmt.Fprintf(pemFile, "-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----\n", base64Der); err != nil {
		return err
	}

	// and the der file is just the binary version of it
	var base64Decoder = base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Der))

	derFile, err := os.Create(DERFilename)
	if err != nil {
		return err
	}
	defer derFile.Close()

	if _, err = io.Copy(derFile, base64Decoder); err != nil {
		return err
	}

	// let's check if we have to update entry in the db:
	exists := slices.ContainsFunc(cdb.certs, func(cert CertificateFile) bool {
		return cert.Environment == environment && slices.Equal(cert.Usage, usage)
	})

	if !exists {
		cdb.certs = append(cdb.certs, CertificateFile{
			Environment: environment,
			Usage:       usage,
			PEMFile:     PEMFilename,
			DERFile:     DERFilename,
		})

		cdb.dirty = true
	}

	return nil
}
