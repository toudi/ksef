package client

import (
	"context"
	"ksef/core/certsdb"
	"ksef/core/invoicedb"
	"ksef/core/keyring"
	"ksef/core/ksef"
	v2 "ksef/internal/client/v2"
	"ksef/internal/http"
)

// InvoicePDFRenderer renders a PDF for an already-stored invoice.
type InvoicePDFRenderer interface {
	RenderInvoice(
		ctx context.Context,
		stored *invoicedb.StoredInvoice,
	) error
}

type ClientInitializer func(*Client)

func WithPDFRenderer(renderer InvoicePDFRenderer) ClientInitializer {
	return func(c *Client) {
		c.pdfRenderer = renderer
	}
}

type Client struct {
	certsDB     *certsdb.CertsDB
	invoiceDB   *invoicedb.InvoiceDB
	keyring     keyring.Keyring
	httpClient  *http.Client
	gateway     ksef.KSeFGatewayConfig
	pdfRenderer InvoicePDFRenderer
	apiClient   *v2.APIClient
}

func Init(
	certsDB *certsdb.CertsDB,
	invoiceDB *invoicedb.InvoiceDB,
	keyring keyring.Keyring,
	gateway ksef.KSeFGatewayConfig,
	initializers ...ClientInitializer,
) *Client {
	client := &Client{
		certsDB:    certsDB,
		invoiceDB:  invoiceDB,
		keyring:    keyring,
		httpClient: http.NewClient(gateway.APIBaseURL),
		gateway:    gateway,
	}
	for _, initializer := range initializers {
		initializer(client)
	}
	return client
}
