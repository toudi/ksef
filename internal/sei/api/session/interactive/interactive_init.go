package interactive

import (
	"fmt"
	"ksef/internal/invoice"
	registryPkg "ksef/internal/registry"
	"ksef/internal/sei/api/client"
	"os"
	"path"
)

type InteractiveSession struct {
	token       string
	apiClient   *client.APIClient
	session     *client.HTTPSession
	referenceNo string
	// token obtained from KSeF web app
	issuerToken string
}

func InteractiveSessionInit(apiClient *client.APIClient) *InteractiveSession {
	session := &InteractiveSession{apiClient: apiClient}

	return session
}

var invoicePayload invoicePayloadType

func (i *InteractiveSession) UploadInvoices(sourcePath string) error {
	var err error

	collection, err := invoice.InvoiceCollection(sourcePath)
	if err != nil {
		return fmt.Errorf("cannot parse invoice collection: %v", err)
	}
	if err = i.Login(collection.Issuer, false); err != nil {
		return fmt.Errorf("cannot login to gateway: %v", err)
	}

	// upload files
	for _, file := range collection.Files {
		if err = i.sendInvoice(file, &invoicePayload); err != nil {
			return fmt.Errorf("error sending invoice: %v", err)
		}
	}

	if err = i.Logout(); err != nil {
		return fmt.Errorf("cannot logout: %v", err)
	}

	registry := registryPkg.NewRegistry()
	registry.SessionID = i.referenceNo
	registry.Environment = i.apiClient.EnvironmentAlias
	registry.Issuer = collection.Issuer

	return registry.Save(path.Join(sourcePath, "registry.yaml"))
}

// SetIssuerToken populates the issuer token from plaintext
// this is a fallback mechanism for people that cannot use org.freedesktop.secret
// service
func (i *InteractiveSession) SetIssuerToken(tokenSource string) {
	i.issuerToken = os.Getenv(tokenSource)
	// if the environment variable was empty then maybe this is a token given as
	// verbatim ?
	if i.issuerToken == "" {
		i.issuerToken = tokenSource
		// if the issuerToken will still be empty there's nothing that can be done
	}
}

func (i *InteractiveSession) HTTPSession() *client.HTTPSession {
	if i.session == nil {
		i.session = client.NewHTTPSession(i.apiClient.Environment.Host)
	}
	return i.session
}

func init() {
	invoicePayload.InvoiceHash.HashSHA.Algorithm = "SHA-256"
	invoicePayload.InvoiceHash.HashSHA.Encoding = "Base64"

	invoicePayload.InvoicePayload.Type = "encrypted"
	invoicePayload.InvoicePayload.EncryptedInvoiceHash.HashSHA.Algorithm = "SHA-256"
	invoicePayload.InvoicePayload.EncryptedInvoiceHash.HashSHA.Encoding = "Base64"
}
