package client

import (
	"ksef/lib/certsdb"
	invoiceregistry "ksef/lib/invoice-registry"
)

type Client struct {
	certsDB         certsdb.CertificatesDB
	invoiceRegistry invoiceregistry.InvoiceRegistry
}

func Init(
	certsDB certsdb.CertificatesDB,
	invoiceRegistry invoiceregistry.InvoiceRegistry,
) *Client {
	return &Client{
		certsDB:         certsDB,
		invoiceRegistry: invoiceRegistry,
	}
}
