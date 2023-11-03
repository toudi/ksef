package api

import (
	"fmt"
	"ksef/common"
)

type InteractiveSession struct {
	token string
	api   *API
}

func (a *API) InteractiveSessionInit() *InteractiveSession {
	session := &InteractiveSession{api: a}

	return session
}

var invoicePayload invoicePayloadType

func (i *InteractiveSession) UploadInvoices(sourcePath string) error {
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

	return nil
}

func init() {
	invoicePayload.InvoiceHash.HashSHA.Algorithm = "SHA-256"
	invoicePayload.InvoiceHash.HashSHA.Encoding = "Base64"

	invoicePayload.InvoicePayload.Type = "encrypted"
	invoicePayload.InvoicePayload.EncryptedInvoiceHash.HashSHA.Algorithm = "SHA-256"
	invoicePayload.InvoicePayload.EncryptedInvoiceHash.HashSHA.Encoding = "Base64"
}
