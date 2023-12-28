package pdf

import (
	"bytes"
	"fmt"
	"ksef/internal/registry"
	"ksef/internal/utils"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
)

type GotenbergPrinter struct {
	host         string
	templatePath string
}

func (g *GotenbergPrinter) Print(
	contentBase64 string,
	invoiceMeta registry.Invoice,
	output string,
) error {
	templateDir := path.Dir(g.templatePath)
	filesList, err := os.ReadDir(templateDir)
	if err != nil {
		return fmt.Errorf("unable to read files from provided template path: %v", err)
	}

	requestBuffer := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(requestBuffer)

	for _, file := range filesList {
		if file.IsDir() || file.Name() == "index.html" {
			// we want to send out render.html, not the index.html itself
			continue
		}
		fileReader, err := os.Open(path.Join(templateDir, file.Name()))
		if err != nil {
			return fmt.Errorf("cannot open source file %s: %v", file.Name(), err)
		}

		mime := "application/octet-stream"
		filename := file.Name()

		if file.Name() == "render.html" {
			// the main template, which has to be included with the name `index.html`, due to gotenberg's requirements
			filename = "index.html"
			mime = "text/html"
		}

		err = utils.AddMultipartFile(
			multipartWriter,
			"files",
			filename, mime, fileReader,
		)
		fileReader.Close()
		if err != nil {
			return fmt.Errorf("cannot add file to gotenberg: %v", err)
		}

	}

	if err = multipartWriter.Close(); err != nil {
		return fmt.Errorf("unable to close multipartWriter: %v", err)
	}

	endpoint, err := url.Parse(g.host)
	if err != nil {
		return fmt.Errorf("invalid host")
	}
	endpoint.Path = "/forms/chromium/convert/html"

	request, err := http.NewRequest("POST", endpoint.String(), requestBuffer)
	if err != nil {
		return fmt.Errorf("unable to create HTTP request to gotenberg: %v", err)
	}

	request.Header.Add("Content-Type", multipartWriter.FormDataContentType())

	return utils.DownloadFile(request, output)
}
