package uploader

import (
	"bytes"
	"errors"
	"ksef/internal/certsdb"
	"ksef/internal/config"
	"ksef/internal/registry"
	"ksef/internal/runtime"
	"ksef/internal/sei"
	"os"

	"github.com/spf13/viper"
)

type Uploader struct {
	dataDir       string
	prefixDir     string
	invoiceDB     *InvoiceDB
	registry      *registry.InvoiceRegistry
	generator     *sei.SEI
	contentBuffer bytes.Buffer
	offlineCert   *certsdb.Certificate
	certsDB       *certsdb.CertificatesDB
}

func New(vip *viper.Viper) (*Uploader, error) {
	dataDir := config.DataDir(vip)
	if err := os.MkdirAll(dataDir, 0775); err != nil {
		return nil, err
	}
	certsDB, err := certsdb.OpenOrCreate(runtime.GetGateway(vip))
	if err != nil {
		return nil, err
	}

	return &Uploader{
		dataDir: dataDir,
		certsDB: certsDB,
	}, nil
}

func (u *Uploader) SetGenerator(g *sei.SEI) {
	u.generator = g
}

func (u *Uploader) GetOfflineModeCertificate(issuerNIP string) (*certsdb.Certificate, error) {
	if u.offlineCert == nil {
		offlineCert, err := u.certsDB.GetByUsage(
			certsdb.UsageOffline, issuerNIP,
		)
		if err != nil {
			return nil, err
		}
		u.offlineCert = &offlineCert
	}

	return u.offlineCert, nil
}

func (u *Uploader) Close() error {
	if err := u.invoiceDB.Save(); err != nil {
		return errors.Join(errors.New("error saving invoiceDB"), err)
	}
	return u.registry.Save("")
}
