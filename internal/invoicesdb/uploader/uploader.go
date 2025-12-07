package uploader

import (
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/spf13/viper"
)

type Uploader struct {
	Queue                     map[monthlyregistry.InvoiceFormCode][]string
	registryByInvoiceChecksum map[string]*monthlyregistry.Registry
}

func NewUploader(vip *viper.Viper) *Uploader {
	return &Uploader{
		Queue:                     make(map[monthlyregistry.InvoiceFormCode][]string),
		registryByInvoiceChecksum: make(map[string]*monthlyregistry.Registry),
	}
}
