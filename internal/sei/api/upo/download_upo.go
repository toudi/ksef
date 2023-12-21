package upo

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	registryPkg "ksef/internal/registry"
	"ksef/internal/sei/api/client"
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
	InvoiceNumber          string `xml:"NumerFaktury" json:"invoiceNumber" yaml:"invoiceNumber"`
	KSeFInvoiceReferenceNo string `xml:"NumerKSeFDokumentu" json:"ksefDocumentId" yaml:"ksefDocumentId"`
	DocumentChecksum       string `xml:"SkrotDokumentu"`
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

const endpointStatus = "common/Status/%s"
const qrcodeUrl = "https://%s/web/verify/%s/%s"

func DownloadUPO(a *client.APIClient, registry *registryPkg.InvoiceRegistry, params *DownloadUPOParams) error {
	var upoStatus upoStatusType
	session := client.NewRequestFactory(a)

	_, err := session.JSONRequest("GET", fmt.Sprintf(endpointStatus, registry.SessionID), nil, &upoStatus)
	if err != nil {
		return fmt.Errorf("get UPO status err=%v", err)
	}

	if upoStatus.ProcessingStatus != 200 {
		return fmt.Errorf("unexpected UPO processing status: %d != 200", upoStatus.ProcessingStatus)
	}

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

	for _, invoiceId := range upo.InvoiceIDS {
		registry.Invoices = append(registry.Invoices, registryPkg.Invoice{
			ReferenceNumber:    invoiceId.InvoiceNumber,
			SEIReferenceNumber: invoiceId.KSeFInvoiceReferenceNo,
			SEIQRCode: fmt.Sprintf(
				qrcodeUrl,
				a.Environment.Host,
				invoiceId.KSeFInvoiceReferenceNo,
				url.QueryEscape(invoiceId.DocumentChecksum),
			),
		})
		qr, err := qrsvg.New(registry.Invoices[len(registry.Invoices)-1].SEIQRCode)
		if err == nil {
			// if there's an error outputting the qrcode there's nothing we can do
			// about it anyway.
			_ = os.WriteFile(path.Join(params.OutputPath, invoiceId.KSeFInvoiceReferenceNo+".svg"), []byte(qr.String()), 0644)
		}
	}

	if err = registry.Save(""); err != nil {
		return fmt.Errorf("error saving status info: %v", err)
	}

	if params.OutputFormat == UPOFormatXML {
		return os.WriteFile(params.Output, upoXMLBytes, 0644)
	}

	// otherwise, we have to send a special request:
	upoPDFURL, err := url.Parse(fmt.Sprintf("https://%s/web/api/session/get-upo-pdf-file", a.Environment.Host))
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
