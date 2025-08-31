package pdf

import (
	"bytes"
	"errors"
	"fmt"
	registryPkg "ksef/internal/registry"
	"ksef/internal/sei/api/client"
	"os"
	"path"
)

type DownloadPDFArgs struct {
	Output      string
	IssuerToken string
	Token       string
	SaveXml     bool
	Invoice     string
}

const DownloadInvoicePDF = "/web/api/invoice/get-invoice-pdf-file?ksefReferenceNumber=%s"

var IsNotXMLInvoice error = errors.New("this is not an XML invoice")

func DownloadPDFFromLocalFile(
	apiClient *client.APIClient,
	registry *registryPkg.InvoiceRegistry,
	args *DownloadPDFArgs,
) error {
	// let's check if the specified `invoice` is actually a source XML file.
	invoiceStruct, err := registryPkg.ParseInvoice(args.Invoice)
	if err != nil {
		// the code using this will catch this error and continue via different path
		return IsNotXMLInvoice
	}
	// yes, it is! let's download the PDF based on that.
	seiRefNo, err := registry.GetSEIRefNo(invoiceStruct.InvoiceNumber)
	if err != nil {
		return fmt.Errorf("unable to find the invoice in status file. was it uploaded?")
	}
	sourceInvoiceBytes, err := os.ReadFile(args.Invoice)
	if err != nil {
		return fmt.Errorf("unable to read the source file: %v", err)
	}

	httpSession := client.NewHTTPSession(apiClient.Environment.Host)
	invoiceXMLReader := bytes.NewReader(sourceInvoiceBytes)

	return httpSession.DownloadPDFFromSourceXML(
		fmt.Sprintf(DownloadInvoicePDF, seiRefNo),
		seiRefNo+".xml",
		invoiceXMLReader,
		path.Join(args.Output, seiRefNo+".pdf"),
	)

}
