package upo

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"ksef/internal/logging"
	registryPkg "ksef/internal/registry"
	"ksef/internal/sei/api/client"
	"log/slog"
	"net/url"
	"os"
	"path"

	qrsvg "github.com/wamuir/svg-qr-code"
)

const (
	UPOFormatXML string = "xml"
	UPOFormatPDF string = "pdf"
)

type upoStatusType struct {
	ProcessingStatus uint16 `json:"processingCode"`
	UPOBase64        string `json:"upo"`
}

type KsefInvoiceIdType struct {
	InvoiceNumber          string `xml:"NumerFaktury"       json:"invoiceNumber"  yaml:"invoiceNumber"`
	KSeFInvoiceReferenceNo string `xml:"NumerKSeFDokumentu" json:"ksefDocumentId" yaml:"ksefDocumentId"`
	DocumentChecksum       string `xml:"SkrotDokumentu"`
}

func (iid KsefInvoiceIdType) Checksum() (string, error) {
	// the checksum in upstream type is encoded with base64. we have to decode it
	// to get the array of bytes and then encode it to hex
	checksumBytes, err := base64.StdEncoding.DecodeString(iid.DocumentChecksum)
	if err != nil {
		return "", fmt.Errorf("unable to decode checksum from base64: %v", err)
	}

	return hex.EncodeToString(checksumBytes), nil
}

type UPO struct {
	XMLName    xml.Name            `xml:"Potwierdzenie"`
	InvoiceIDS []KsefInvoiceIdType `xml:"Dokument"`
}

type DownloadUPOParams struct {
	OutputPath   string
	Output       string
	Mkdir        bool
	OutputFormat string
}

const endpointStatus = "/api/common/Status/%s"
const qrcodeUrl = "https://%s/web/verify/%s/%s"

func DownloadUPO(
	a *client.APIClient,
	registry *registryPkg.InvoiceRegistry,
	params *DownloadUPOParams,
) error {
	var log *slog.Logger = logging.UPOLogger
	log.Info("DownloadUPO", "sessionId", registry.SessionID)

	var upoStatus upoStatusType
	session := client.NewHTTPSession(a.Environment.Host)

	_, err := session.JSONRequest(
		client.JSONRequestParams{
			Method:   "GET",
			Endpoint: fmt.Sprintf(endpointStatus, registry.SessionID),
			Payload:  nil,
			Response: &upoStatus,
			Logger:   logging.UPOHTTPLogger,
		},
	)
	if err != nil {
		return fmt.Errorf("get UPO status err=%v", err)
	}

	if upoStatus.ProcessingStatus != 200 {
		return fmt.Errorf("unexpected UPO processing status: %d != 200", upoStatus.ProcessingStatus)
	}
	log.Debug("UPO is ready")

	// we have to decode UPO into xml regardless of what we decide to do next

	upoXMLBytes, err := base64.StdEncoding.DecodeString(upoStatus.UPOBase64)
	if err != nil {
		return fmt.Errorf("unable to decode UPO XML from base64: %v", err)
	}

	var upo UPO
	// parse upo to obtain KSeF reference numbers for the invoice numbers
	if err = xml.Unmarshal(upoXMLBytes, &upo); err != nil {
		return fmt.Errorf("unable to parse upo as XML structure: %v", err)
	}
	log.Debug("Parsed UPO XML")

	for _, invoiceId := range upo.InvoiceIDS {
		log.Info(
			"processing",
			"invoice id",
			invoiceId.InvoiceNumber,
			"ksef invoice id",
			invoiceId.KSeFInvoiceReferenceNo,
		)

		checksum, err := invoiceId.Checksum()
		if err != nil {
			return fmt.Errorf("unable to obtain upstream document checksum: %v", err)
		}

		invoice, err := registry.Update(registryPkg.Invoice{
			ReferenceNumber:    invoiceId.InvoiceNumber,
			SEIReferenceNumber: invoiceId.KSeFInvoiceReferenceNo,
			SEIQRCode: fmt.Sprintf(
				qrcodeUrl,
				a.Environment.Host,
				invoiceId.KSeFInvoiceReferenceNo,
				url.QueryEscape(invoiceId.DocumentChecksum),
			),
			Checksum: checksum,
		})

		if err == registryPkg.ErrDoesNotExist {
			return fmt.Errorf("UPO contains an invoice which does not exist in registry")
		}

		if err != nil {
			return fmt.Errorf("unexpected error while updating registry: %v", err)
		}

		qr, err := qrsvg.New(invoice.SEIQRCode)

		if err == nil {
			// if there's an error outputting the qrcode there's nothing we can do
			// about it anyway.
			_ = os.WriteFile(
				path.Join(params.OutputPath, invoiceId.KSeFInvoiceReferenceNo+".svg"),
				[]byte(qr.String()),
				0644,
			)
		}
	}

	if err = registry.Save(""); err != nil {
		return fmt.Errorf("error saving status info: %v", err)
	}

	if params.OutputFormat == UPOFormatXML {
		log.Debug("save upo as XML file")
		return os.WriteFile(params.Output, upoXMLBytes, 0644)
	}

	// otherwise, we have to send a special request:
	log.Info("fetch UPO as PDF")
	upoPDFURL, err := url.Parse(
		fmt.Sprintf("https://%s/web/api/session/get-upo-pdf-file", a.Environment.Host),
	)
	if err != nil {
		return fmt.Errorf("unable to parse url for UPO PDF")
	}

	return session.DownloadPDFFromSourceXML(
		upoPDFURL.String(),
		registry.SessionID,
		bytes.NewReader(upoXMLBytes),
		params.Output,
	)
}
