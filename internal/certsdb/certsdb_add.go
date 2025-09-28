package certsdb

import (
	"encoding/base64"
	"fmt"
	"io"
	"ksef/internal/environment"
	"os"
	"path"
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

	PEMFilename := path.Dir(certificatesDBFile) + certBaseFilename + ".pem"
	DERFilename := path.Dir(certificatesDBFile) + certBaseFilename + ".der"
	// so the base64-encoded content is actually the essense of PEM so let's use a nifty hack to save it
	pemFile, err := os.Create(PEMFilename)
	if err != nil {
		return err
	}
	defer pemFile.Close()
	if _, err = fmt.Fprintf(pemFile, "-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----\n", base64Der); err != nil {
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

	cdb.certs = append(cdb.certs, CertificateFile{
		Environment: environment,
		Usage:       usage,
		PEMFile:     PEMFilename,
		DERFile:     DERFilename,
	})

	cdb.dirty = true

	return nil
}
