package interactive

import (
	"errors"
	"fmt"
	"ksef/internal/invoice"
	"ksef/internal/logging"
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

type interactiveSessionUploadParams struct {
	ForceUpload bool
}

func InteractiveSessionInit(apiClient *client.APIClient) *InteractiveSession {
	session := &InteractiveSession{apiClient: apiClient}

	return session
}

var invoicePayload invoicePayloadType
var InteractiveSessionUploadParams interactiveSessionUploadParams
var ErrProbablyUsedSend = errors.New("upload command probably used previously")

func (i *InteractiveSession) UploadInvoices(sourcePath string) error {
	var err error

	// we try to load the registry. If it does not exists, it's not a fatal problem since
	// we will receive a new registry as the return value and it will still conform to the
	// interface, it will simply report each invoice as not sent.
	registry, err := registryPkg.OpenOrCreate(path.Join(sourcePath, "registry.yaml"))
	if err != nil {
		return fmt.Errorf("cannot open registry: %v", err)
	}

	collection, err := invoice.InvoiceCollection(sourcePath, registry)

	if err == invoice.ErrAlreadySynced {
		logging.UploadLogger.Info("no invoices left to send")
		return nil
	}

	// ok, there are some potential invoices, however let's warn the user if the registry already contains
	// the session ID as well
	if registry.SessionID != "" && !InteractiveSessionUploadParams.ForceUpload {
		return ErrProbablyUsedSend
	}

	if err != nil {
		return fmt.Errorf("cannot parse invoice collection: %v", err)
	}

	// check NIP for validity
	if !i.apiClient.Environment.NipValidator(collection.Issuer) {
		return fmt.Errorf("invalid NIP: %s", collection.Issuer)
	}

	if err = i.Login(collection.Issuer, false); err != nil {
		return fmt.Errorf("cannot login to gateway: %v", err)
	}

	// upload files
	for _, file := range collection.Files {
		if err = i.sendInvoice(file.Filename, &invoicePayload); err != nil {
			return fmt.Errorf("error sending invoice: %v", err)
		}
	}

	if err = i.Logout(); err != nil {
		return fmt.Errorf("cannot logout: %v", err)
	}

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
