package client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strings"
)

func (hs *HTTPSession) DownloadPDFFromSourceXML(
	endpoint string,
	sourceXMLFileName string,
	sourceXMLFile io.Reader,
	outputPath string,
) error {
	if !strings.HasPrefix(endpoint, "https") {
		endpoint = fmt.Sprintf("https://%s%s", hs.host, endpoint)
	}
	pdfURL, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("unable to parse url for PDF")
	}

	requestBuffer := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(requestBuffer)
	xmlHeader := make(textproto.MIMEHeader)
	xmlHeader.Set(
		"Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s.xml"`, "file", sourceXMLFileName),
	)
	xmlHeader.Set("Content-Type", "text/xml")
	xmlFile, err := multipartWriter.CreatePart(xmlHeader)
	if err != nil {
		return fmt.Errorf("unable to create xml file writer: %v", err)
	}
	if _, err = io.Copy(xmlFile, sourceXMLFile); err != nil {
		return fmt.Errorf("unable to write xml bytes to HTTP request: %v", err)
	}
	if err = multipartWriter.Close(); err != nil {
		return fmt.Errorf("unable to close multipartWriter: %v", err)
	}

	pdfDownloadRequest, err := http.NewRequest("POST", pdfURL.String(), requestBuffer)
	if err != nil {
		return fmt.Errorf("unable to prepare pdf download request: %v", err)
	}
	pdfDownloadRequest.Header.Add("Content-Type", multipartWriter.FormDataContentType())

	response, err := http.DefaultClient.Do(pdfDownloadRequest)
	if err != nil {
		return fmt.Errorf("unable to perform download request: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected response from pdf download: %d != 200", response.StatusCode)
	}

	pdfFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("unable to create PDF file: %v", err)
	}
	defer pdfFile.Close()

	_, err = io.Copy(pdfFile, response.Body)

	return err
}
