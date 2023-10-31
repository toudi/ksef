package uploader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Session struct {
	ReferenceNumber string `json:"referenceNumber"`
	Token           struct {
		Value string `json:"token"`
	} `json:"sessionToken"`
}

func (u *Uploader) uploadInteractive(sourcePath string) error {
	var err error
	var session *Session

	if session, err = u.initSession(); err != nil {
		return fmt.Errorf("błąd inicjalizacji sesji: %v", err)
	}
	// pora na faktyczny upload
	var invoicePayload invoicePayloadType
	invoicePayload.InvoiceHash.HashSHA.Algorithm = "SHA-256"
	invoicePayload.InvoiceHash.HashSHA.Encoding = "Base64"

	invoicePayload.InvoicePayload.Type = "encrypted"
	invoicePayload.InvoicePayload.EncryptedInvoiceHash.HashSHA.Algorithm = "SHA-256"
	invoicePayload.InvoicePayload.EncryptedInvoiceHash.HashSHA.Encoding = "Base64"

	files, err := os.ReadDir(sourcePath)
	if err != nil {
		return fmt.Errorf("cannot read list of files from %s: %v", sourcePath, err)
	}

	var fileName string

	for _, file := range files {
		fileName = file.Name()

		if filepath.Ext(fileName) != ".xml" || filepath.Base(fileName) == metadataFileName {
			continue
		}

		err = u.interactiveSendFile(sourcePath+string(os.PathSeparator)+fileName, &invoicePayload, session)
		if err != nil {
			return fmt.Errorf("błąd przesyłania faktury: %v", err)
		}
	}

	// uzywamy numeru sesji zeby pobrac UPO
	fmt.Printf("session: %s\n", session.Token.Value)

	// zakonczmy sesje
	terminateRequest, err := http.NewRequest("GET", u.host+"api/online/Session/Terminate", nil)
	if err != nil {
		return fmt.Errorf("błąd tworzenia requestu Session/Terminate: %v", err)
	}
	terminateRequest.Header.Set("SessionToken", session.Token.Value)
	terminateResponse, err := http.DefaultClient.Do(terminateRequest)
	if err != nil || terminateResponse.StatusCode/100 != 2 {
		defer terminateResponse.Body.Close()
		terminateResponseBody, _ := io.ReadAll(terminateResponse.Body)
		return fmt.Errorf("błąd kończenia sesji: statuscode=%d, err=%v\n%s", terminateResponse.StatusCode, err, string(terminateResponseBody))
	}

	return nil
}
