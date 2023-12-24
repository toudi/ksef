package interactive

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"ksef/internal/logging"
	"ksef/internal/sei/api/client"
	"log/slog"
	"os"
)

type fileHashMetadata struct {
	HashSHA struct {
		Algorithm string `json:"algorithm"`
		Encoding  string `json:"encoding"`
		Value     string `json:"value"`
	} `json:"hashSHA"`
	FileSize int `json:"fileSize"`
}

type invoicePayloadType struct {
	InvoiceHash    fileHashMetadata `json:"invoiceHash"`
	InvoicePayload struct {
		Type                 string           `json:"type"`
		EncryptedInvoiceHash fileHashMetadata `json:"encryptedInvoiceHash"`
		EncryptedInvoiceBody string           `json:"encryptedInvoiceBody"`
	} `json:"invoicePayload"`
}

const endpointSendInvoice = "/api/online/Invoice/Send"

func (i *InteractiveSession) sendInvoice(
	filePath string,
	invoicePayload *invoicePayloadType,
) error {
	var log *slog.Logger = logging.UploadLogger
	log.Info("InteractiveSession::sendInvoice", "source file", filePath)

	invoiceContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("unable to read invoice content: %v", err)
	}

	hasher := sha256.New()
	hasher.Write(invoiceContent)

	log.Debug("InteractiveSession::sendInvoice - hash file")
	invoicePayload.InvoiceHash.FileSize = len(invoiceContent)
	invoicePayload.InvoiceHash.HashSHA.Value = base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	encryption, _ := i.apiClient.Encryption()

	// encrypt file
	log.Debug("InteractiveSession::sendInvoice - encrypt file")
	encryptedContent := encryption.Cipher.Encrypt(invoiceContent, true)

	hasher.Reset()
	hasher.Write(encryptedContent)

	invoicePayload.InvoicePayload.EncryptedInvoiceBody = base64.StdEncoding.EncodeToString(
		encryptedContent,
	)
	invoicePayload.InvoicePayload.EncryptedInvoiceHash.FileSize = len(encryptedContent)
	invoicePayload.InvoicePayload.EncryptedInvoiceHash.HashSHA.Value = base64.StdEncoding.EncodeToString(
		hasher.Sum(nil),
	)

	invoiceSendResponse, err := i.session.JSONRequest(client.JSONRequestParams{
		Method:   "PUT",
		Endpoint: endpointSendInvoice,
		Payload:  &invoicePayload,
		Response: nil,
		Logger:   logging.UploadHTTPLogger,
	})
	if err != nil {
		return fmt.Errorf("error sending invoice: %v", err)
	}
	if invoiceSendResponse.StatusCode/100 != 2 {
		return fmt.Errorf("unexpected response code: %d != 2xx", invoiceSendResponse.StatusCode)
	}

	return nil
}
