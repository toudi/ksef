package api

import (
	"fmt"
	"ksef/common"
	"os"
)

type InteractiveSession struct {
	token       string
	api         *API
	referenceNo string
	// token obtained from KSeF web app
	issuerToken string
}

func (a *API) InteractiveSessionInit() *InteractiveSession {
	session := &InteractiveSession{api: a}

	return session
}

var invoicePayload invoicePayloadType

func (i *InteractiveSession) UploadInvoices(sourcePath string, statusFileFormat string) error {
	var err error

	collection, err := common.InvoiceCollection(sourcePath)
	if err != nil {
		return fmt.Errorf("cannot parse invoice collection: %v", err)
	}
	if err = i.login(collection.Issuer); err != nil {
		return fmt.Errorf("cannot login to gateway: %v", err)
	}

	// upload files
	for _, file := range collection.Files {
		if err = i.sendInvoice(file, &invoicePayload); err != nil {
			return fmt.Errorf("error sending invoice: %v", err)
		}
	}

	if err = i.logout(); err != nil {
		return fmt.Errorf("cannot logout: %v", err)
	}

	i.token = ""
	i.api.cipher = nil

	return (&StatusInfo{
		selectedFormat: statusFileFormat,
		SessionID:      i.referenceNo,
		Environment:    i.api.environmentAlias,
	}).Save(sourcePath)
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

func init() {
	invoicePayload.InvoiceHash.HashSHA.Algorithm = "SHA-256"
	invoicePayload.InvoiceHash.HashSHA.Encoding = "Base64"

	invoicePayload.InvoicePayload.Type = "encrypted"
	invoicePayload.InvoicePayload.EncryptedInvoiceHash.HashSHA.Algorithm = "SHA-256"
	invoicePayload.InvoicePayload.EncryptedInvoiceHash.HashSHA.Encoding = "Base64"
}
