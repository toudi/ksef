package client

import (
	"bytes"
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"text/template"
)

type RequestFactory struct {
	api     *APIClient
	headers http.Header
}

func NewRequestFactory(api *APIClient) *RequestFactory {
	return &RequestFactory{api: api, headers: make(http.Header)}
}

func (rf *RequestFactory) SetHeader(header string, value string) {
	rf.headers.Add(header, value)
}

func (rf *RequestFactory) Request(method string, endpoint string, payload io.Reader) (*http.Request, error) {
	// if it's not a public URL already ..
	if !strings.HasPrefix(endpoint, "http") {
		// convert it to one
		endpoint = rf.api.apiEndpoint(endpoint)
	}
	request, err := http.NewRequest(method, endpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	for key, values := range rf.headers {
		request.Header.Add(key, values[0])
	}

	return request, nil
}

func (rf *RequestFactory) JSONRequest(method string, endpoint string, payload interface{}, response interface{}) (*http.Response, error) {
	var encodedPayload bytes.Buffer
	if err := json.NewEncoder(&encodedPayload).Encode(payload); err != nil {
		return nil, fmt.Errorf("error encoding JSON: %v", err)
	}
	request, err := rf.Request(method, endpoint, &encodedPayload)
	// fmt.Printf("request URL: %s\n", request.URL.String())
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")

	httpResponse, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error performing HTTP request: %v", err)
	}
	defer httpResponse.Body.Close()
	if response != nil {
		httpResponseBody, _ := io.ReadAll(httpResponse.Body)
		// fmt.Printf("response body: \n%s\n", string(httpResponseBody))
		err = json.NewDecoder(bytes.NewReader(httpResponseBody)).Decode(response)
		if err != nil {
			return nil, fmt.Errorf("error decoding response: %v", err)
		}
	}
	return httpResponse, nil
}

func (rf *RequestFactory) XMLRequest(method string, endpoint string, templateDir embed.FS, templateName string, templateData interface{}, response interface{}) (*http.Response, error) {
	var funcMap = template.FuncMap{
		"base64": base64.StdEncoding.EncodeToString,
	}

	tmpl, err := template.New(templateName).Funcs(funcMap).ParseFS(templateDir, templateName)

	if err != nil {
		return nil, fmt.Errorf("error initialising template: %v", err)
	}

	var renderedTemplate bytes.Buffer

	if err = tmpl.Execute(
		&renderedTemplate,
		templateData,
	); err != nil {
		return nil, fmt.Errorf("error rendering authRequest template: %v", err)
	}

	// fmt.Printf("posting xml template: \n%s\n", renderedTemplate.String())

	request, err := rf.Request(method, endpoint, &renderedTemplate)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}
	request.Header.Set("Content-Type", "application/octet-stream")

	httpResponse, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error performing HTTP request: %v", err)
	}
	defer httpResponse.Body.Close()
	if response != nil {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		// fmt.Printf("response body: \n%s\n", string(responseBody))
		err = json.NewDecoder(bytes.NewReader(responseBody)).Decode(response)
		if err != nil {
			return nil, fmt.Errorf("error decoding response: %v", err)
		}
	}
	return httpResponse, nil
}

func (rf *RequestFactory) SendFile(method string, endpoint string, fileName string, jsonResponse interface{}) (*http.Response, error) {
	fileReader, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("unable to open file for reading: %v", err)
	}
	defer fileReader.Close()

	request, err := rf.Request(method, endpoint, fileReader)

	if err != nil {
		return nil, fmt.Errorf("error preparing request: %v", err)
	}
	// request.Header.Add("Content-Type", "application/octet-stream")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error sending file: %v", err)
	}
	defer response.Body.Close()
	fmt.Printf("initResponse status: %d\n", response.StatusCode)
	if response.StatusCode/100 != 2 {
		responseContent, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response:%v", err)
		}
		return nil, fmt.Errorf("unexpected response code from initResponse:\n%s\n", string(responseContent))
	}

	err = json.NewDecoder(response.Body).Decode(jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %v", err)
	}

	return response, nil
}

func (rf *RequestFactory) DownloadPDFFromSourceXML(endpoint string, sourceXMLFileName string, sourceXMLFile io.Reader, outputPath string) error {
	pdfURL, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("unable to parse url for PDF")
	}

	requestBuffer := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(requestBuffer)
	xmlHeader := make(textproto.MIMEHeader)
	xmlHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s.xml"`, "file", sourceXMLFileName))
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
