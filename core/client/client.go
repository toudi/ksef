package client

import (
	"ksef/core/certsdb"
	"ksef/core/invoicedb"
	"ksef/core/keyring"
	"ksef/internal/http"
)

type Client struct {
	certsDB    *certsdb.CertsDB
	invoiceDB  *invoicedb.InvoiceDB
	keyring    keyring.Keyring
	httpClient *http.Client
}

func Init(
	certsDB *certsdb.CertsDB,
	invoiceDB *invoicedb.InvoiceDB,
	keyring keyring.Keyring,
	httpClient *http.Client,
) *Client {
	return &Client{
		certsDB:    certsDB,
		invoiceDB:  invoiceDB,
		keyring:    keyring,
		httpClient: httpClient,
	}
}
