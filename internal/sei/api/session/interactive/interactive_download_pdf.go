package interactive

import (
	"bytes"
	"fmt"
	"io"
	"ksef/internal/logging"
	"ksef/internal/pdf"
	"ksef/internal/registry"
	"ksef/internal/sei/api/client"
	"net/http"
	"os"
	"path"
)

const downloadInvoiceXML = "/api/online/Invoice/Get/%s"

func DownloadPDFFromAPI(
	apiClient *client.APIClient,
	args *pdf.DownloadPDFArgs,
	r *registry.InvoiceRegistry,
) error {
	// we have to initialize the interactive session
	session := InteractiveSessionInit(apiClient)
	httpSession := session.HTTPSession()
	var err error

	// do we have a token?
	if args.Token != "" {
		// yes! this means we can skip the lengthy login procedure
		httpSession.SetHeader("SessionToken", args.Token)
	} else {
		// no luck this time. let's proceed with login
		if args.IssuerToken != "" {
			session.SetIssuerToken(args.IssuerToken)
		}
		if err := session.Login(r.Issuer, true); err != nil {
			return fmt.Errorf("unable to login to interactive session: %v", err)
		}
	}
	// defer session.Logout()

	// httpSession := session.HTTPSession()
	var seiRefNo string
	var downloadAll = args.Invoice == "*"

	for _, invoice := range r.Invoices {
		// * means to download all invoices.
		seiRefNo = invoice.SEIReferenceNumber
		if args.Invoice != "*" {
			if seiRefNo, err = r.GetSEIRefNo(args.Invoice); err != nil {
				return fmt.Errorf("cannot find invoice %s in registry", args.Invoice)
			}
		}
		if seiRefNo == "" {
			// we don't have the sei ref no for this ivnoice.
			continue
		}
		if _, err = os.Stat(path.Join(args.Output, seiRefNo+".pdf")); err == nil {
			// PDF was already downloaded
			if downloadAll {
				// nothing left to do
				return nil
			}
			continue
		}

		var invoiceXMLReader io.Reader
		var xmlFilename = path.Join(args.Output, seiRefNo+".xml")

		if _, err = os.Stat(xmlFilename); err != nil {
			invoiceXMLRequest, err := httpSession.Request(
				"GET",
				fmt.Sprintf(downloadInvoiceXML, seiRefNo),
				nil,
				logging.InteractiveHTTPLogger,
			)
			if err != nil {
				return fmt.Errorf("unable to download invoice in XML Format: %v", err)
			}
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
			if args.SaveXml {
				if err = os.WriteFile(xmlFilename, invoiceXML.Bytes(), 0644); err != nil {
					return fmt.Errorf("unable to persist XML file to disk: %v", err)
				}
			}
			invoiceXMLReader = bytes.NewReader(invoiceXML.Bytes())
		} else {
			if invoiceXMLReader, err = os.Open(xmlFilename); err != nil {
				return fmt.Errorf("unable to read XML file: %v", err)
			}
		}

		err = httpSession.DownloadPDFFromSourceXML(
			fmt.Sprintf(pdf.DownloadInvoicePDF, seiRefNo),
			seiRefNo+".xml",
			invoiceXMLReader,
			path.Join(args.Output, seiRefNo+".pdf"),
		)

		if err != nil {
			return fmt.Errorf("unable to download PDF: %v", err)
		}
	}

	return nil
}
