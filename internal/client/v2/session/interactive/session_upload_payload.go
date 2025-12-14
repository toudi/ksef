package interactive

import (
	"encoding/base64"
	"ksef/internal/utils"
	"os"
)

type SessionUploadPayload struct {
	InvoiceHash             string `json:"invoiceHash"`
	InvoiceSize             int    `json:"invoiceSize"`
	EncryptedInvoiceHash    string `json:"encryptedInvoiceHash"`
	EncryptedInvoiceSize    int    `json:"encryptedInvoiceSize"`
	EncryptedInvoiceContent string `json:"encryptedInvoiceContent"`
	OfflineMode             bool   `json:"offlineMode"`
}

func (s *Session) getUploadPayload(us *uploadSession, filename string, offline bool) (*SessionUploadPayload, error) {
	var payload SessionUploadPayload
	var err error

	// first of all, let's open the invoice file
	invoiceBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// now we can start populating the payload.
	payload.InvoiceSize = len(invoiceBytes)
	if payload.InvoiceHash, err = utils.Sha256Base64(invoiceBytes); err != nil {
		return nil, err
	}
	// now for the encryption. this isn't as bad as it seems:
	encryptedInvoiceBytes := us.cipher.Encrypt(invoiceBytes, true)
	payload.EncryptedInvoiceSize = len(encryptedInvoiceBytes)
	payload.EncryptedInvoiceContent = base64.StdEncoding.EncodeToString(encryptedInvoiceBytes)
	if payload.EncryptedInvoiceHash, err = utils.Sha256Base64(encryptedInvoiceBytes); err != nil {
		return nil, err
	}
	payload.OfflineMode = offline

	return &payload, nil
}
