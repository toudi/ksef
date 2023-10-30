package uploader

import (
	"fmt"
	"net/http"
)

type Session struct {
	ReferenceNumber string `json:"referenceNumber"`
}

type fileHashMetadata struct {
	hashSHA struct {
		algorithm string
		encoding  string
		value     string
	}
	fileSize int
}

type invoicePayloadType struct {
	invoiceHash    fileHashMetadata
	invoicePayload struct {
		_type                string
		encryptedInvoiceHash fileHashMetadata
		encryptedInvoiceBody string
	}
}

func (u *Uploader) uploadInteractive(sourcePath string) error {
	var err error
	var session *Session

	if session, err = u.initSession(); err != nil {
		return fmt.Errorf("błąd inicjalizacji sesji: %v", err)
	}
	// pora na faktyczny upload
	var invoicePayload invoicePayloadType
	invoicePayload.invoiceHash.hashSHA.algorithm = "SHA-256"
	invoicePayload.invoiceHash.hashSHA.encoding = "Base64"

	invoicePayload.invoicePayload.encryptedInvoiceHash.hashSHA.algorithm = "SHA-256"
	invoicePayload.invoicePayload.encryptedInvoiceHash.hashSHA.encoding = "Base64"

	// for file in files {

	// }

	// zakonczmy sesje
	terminateResponse, err := http.DefaultClient.Get(u.host + "api/online/Session/Terminate")
	if err != nil || terminateResponse.StatusCode/100 != 2 {
		return fmt.Errorf("błąd kończenia sesji: statuscode=%d, err=%v", terminateResponse.StatusCode, err)
	}

	return nil
}
