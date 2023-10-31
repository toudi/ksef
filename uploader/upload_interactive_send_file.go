package uploader

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func (u *Uploader) interactiveSendFile(filePath string, invoicePayload *invoicePayloadType, session *Session) error {
	invoiceContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("unable to read invoice content: %v", err)
	}

	hasher := sha256.New()
	hasher.Write(invoiceContent)

	invoicePayload.InvoiceHash.FileSize = len(invoiceContent)
	invoicePayload.InvoiceHash.HashSHA.Value = base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	// encrypt file
	encryptedContent := u.cipher.Encrypt(invoiceContent, true)

	hasher.Reset()
	hasher.Write(encryptedContent)

	invoicePayload.InvoicePayload.EncryptedInvoiceBody = base64.StdEncoding.EncodeToString(encryptedContent)
	invoicePayload.InvoicePayload.EncryptedInvoiceHash.FileSize = len(encryptedContent)
	invoicePayload.InvoicePayload.EncryptedInvoiceHash.HashSHA.Value = base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	var invoiceSendPayload bytes.Buffer
	err = json.NewEncoder(&invoiceSendPayload).Encode(invoicePayload)
	if err != nil {
		return fmt.Errorf("błąd wysyłki faktury: %v", err)
	}

	invoiceSendRequest, err := http.NewRequest("PUT", u.host+"api/online/Invoice/Send", &invoiceSendPayload)
	if err != nil {
		return fmt.Errorf("błąd tworzenia requestu Invoice/Send: %v", err)
	}
	invoiceSendRequest.Header.Set("SessionToken", session.Token.Value)
	invoiceSendRequest.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(invoiceSendRequest)

	if err != nil || response.StatusCode/100 != 2 {
		defer response.Body.Close()
		responseBody, _ := io.ReadAll(response.Body)
		fmt.Printf("send body: %s\n", string(responseBody))

		return fmt.Errorf("błąd przesyłania faktury: %v", err)
	}
	fmt.Printf("Przesyłanie faktury: %d\n", response.StatusCode)

	return nil
}
