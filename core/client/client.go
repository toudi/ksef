package client

import (
	"ksef/core/certsdb"
	"ksef/core/invoicedb"
)

type Client struct {
	certsDB   *certsdb.CertsDB
	invoiceDB *invoicedb.InvoiceDB
}

func Init(
	certsDB *certsdb.CertsDB,
	invoiceDB *invoicedb.InvoiceDB,
) *Client {
	return &Client{
		certsDB:   certsDB,
		invoiceDB: invoiceDB,
	}
}
