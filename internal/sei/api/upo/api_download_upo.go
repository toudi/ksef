package upo

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"ksef/internal/sei/api/client"
	"ksef/internal/sei/api/status"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
)

const (
	UPOFormatXML string = "xml"
	UPOFormatPDF string = "pdf"
)

type upoStatusType struct {
	ProcessingStatus uint16 `json:"processingCode"`
	UPOBase64        string `json:"upo"`
}

type UPO struct {
	XMLName    xml.Name                   `xml:"Potwierdzenie"`
	InvoiceIDS []status.KsefInvoiceIdType `xml:"Dokument"`
}

const endpointStatus = "common/Status/%s"

func DownloadUPO(a *client.APIClient, status *status.StatusInfo, outputFormat string, outputPath string) error {
	var upoStatus upoStatusType
	session := client.NewRequestFactory(a)

	_, err := session.JSONRequest("GET", fmt.Sprintf(endpointStatus, status.SessionID), nil, &upoStatus)
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

	status.InvoiceIds = upo.InvoiceIDS
	if err = status.Save(""); err != nil {
		return fmt.Errorf("error saving status info: %v", err)
	}

	if outputFormat == UPOFormatXML {
		return os.WriteFile(outputPath, upoXMLBytes, 0644)
	}

	// otherwise, we have to send a special request:
	upoPDFURL, err := url.Parse(fmt.Sprintf("https://%s/web/api/session/get-upo-pdf-file", a.Environment.Host))
	if err != nil {
		return fmt.Errorf("unable to parse url for UPO PDF")
	}

	requestBuffer := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(requestBuffer)
	xmlHeader := make(textproto.MIMEHeader)
	xmlHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s.xml"`, "file", status.SessionID))
	xmlHeader.Set("Content-Type", "text/xml")
	xmlFile, err := multipartWriter.CreatePart(xmlHeader)
	if err != nil {
		return fmt.Errorf("unable to create xml file writer: %v", err)
	}
	if _, err = io.Copy(xmlFile, bytes.NewReader(upoXMLBytes)); err != nil {
		return fmt.Errorf("unable to write UPO xml bytes to HTTP request: %v", err)
	}
	if err = multipartWriter.Close(); err != nil {
		return fmt.Errorf("unable to close multipartWriter: %v", err)
	}

	upoPDFDownloadRequest, err := http.NewRequest("POST", upoPDFURL.String(), requestBuffer)
	if err != nil {
		return fmt.Errorf("unable to prepare UPO pdf download request: %v", err)
	}
	upoPDFDownloadRequest.Header.Add("Content-Type", multipartWriter.FormDataContentType())

	response, err := http.DefaultClient.Do(upoPDFDownloadRequest)
	if err != nil {
		return fmt.Errorf("unable to perform upo download request: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		// body, _ := io.ReadAll(response.Body)
		// fmt.Printf("body: %s\n", string(body))
		return fmt.Errorf("unexpected response from upo download: %d != 200", response.StatusCode)
	}

	upoPDFFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("unable to create UPO PDF file: %v", err)
	}
	defer upoPDFFile.Close()

	_, err = io.Copy(upoPDFFile, response.Body)
	return err
}