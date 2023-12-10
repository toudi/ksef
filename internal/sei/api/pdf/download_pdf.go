package pdf

import (
	"bytes"
	"fmt"
	"io"
	"ksef/internal/sei/api/client"
	"ksef/internal/sei/api/status"
	"ksef/internal/sei/api/upload/interactive"
	"net/http"
	"path"
)

const downloadInvoiceXML = "online/Invoice/Get/%s"
const downloadInvoicePDF = "https://%s/web/api/invoice/get-invoice-pdf-file?ksefReferenceNumber=%s"

func DownloadPDF(apiClient *client.APIClient, statusInfo *status.StatusInfo, invoiceNo string, outputPath string) error {
	// we have to initialize the interactive session
	session := interactive.InteractiveSessionInit(apiClient)
	if err := session.Login(statusInfo.Issuer); err != nil {
		return fmt.Errorf("unable to login to interactive session: %v", err)
	}
	defer session.Logout()
	if err := session.WaitForStatusCode(interactive.VerifySecuritySuccess, 15); err != nil {
		return fmt.Errorf("Authorisation was successful however the session is not open for further processing: %v", err)
	}

	httpSession := session.HTTPSession()

	invoiceXMLRequest, err := httpSession.Request("GET", fmt.Sprintf(downloadInvoiceXML, invoiceNo), nil)
	if err != nil {
		return fmt.Errorf("unable to download invoice in XML Format: %v", err)
	}
	// time.Sleep(1 * time.Second)
	fmt.Printf("request: %+v\n", invoiceXMLRequest)
	invoiceXMLResponse, err := http.DefaultClient.Do(invoiceXMLRequest)
	if err != nil {
		return fmt.Errorf("error performing HTTP request: %v", err)
	}
	defer invoiceXMLResponse.Body.Close()

	var invoiceXML bytes.Buffer
	if _, err = io.Copy(&invoiceXML, invoiceXMLResponse.Body); err != nil {
		return fmt.Errorf("unable to save XML invoice to the buffer: %v", err)
	}

	invoiceXMLReader := bytes.NewReader(invoiceXML.Bytes())
	return httpSession.DownloadPDFFromSourceXML(
		fmt.Sprintf(downloadInvoicePDF, apiClient.Environment.Host, invoiceNo),
		invoiceNo+".xml",
		invoiceXMLReader,
		path.Join(outputPath, invoiceNo+".pdf"),
	)
}
