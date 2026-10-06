package client

import (
	"ksef/core/certsdb"
	invoiceregistry "ksef/core/invoice-registry"
)

type Client struct {
	certsDB         *certsdb.CertsDB
	invoiceRegistry invoiceregistry.InvoiceRegistry
}

func Init(
	certsDB *certsdb.CertsDB,
	invoiceRegistry invoiceregistry.InvoiceRegistry,
) *Client {
	return &Client{
		certsDB:         certsDB,
		invoiceRegistry: invoiceRegistry,
	}
}
