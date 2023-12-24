package client

import (
	"bytes"
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"text/template"
)

type XMLRequestParams struct {
	Method       string
	Endpoint     string
	TemplateDir  embed.FS
	TemplateName string
	TemplateData interface{}
	Response     interface{}
	Logger       *slog.Logger
}

func (hs *HTTPSession) XMLRequest(params XMLRequestParams) (*http.Response, error) {
	var funcMap = template.FuncMap{
		"base64": base64.StdEncoding.EncodeToString,
	}
	var log *slog.Logger = params.Logger

	tmpl, err := template.New(params.TemplateName).
		Funcs(funcMap).
		ParseFS(params.TemplateDir, params.TemplateName)

	if err != nil {
		return nil, fmt.Errorf("error initialising template: %v", err)
	}

	var renderedTemplate bytes.Buffer

	if err = tmpl.Execute(
		&renderedTemplate,
		params.TemplateData,
	); err != nil {
		return nil, fmt.Errorf("error rendering authRequest template: %v", err)
	}

	log.Debug("HTTPSEssion::XMLRequest rendered XML", "content", renderedTemplate.String())

	request, err := hs.Request(params.Method, params.Endpoint, &renderedTemplate, log)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}
	request.Header.Set("Content-Type", "application/octet-stream")

	httpResponse, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error performing HTTP request: %v", err)
	}
	defer httpResponse.Body.Close()

	if params.Response != nil {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		log.Debug(
			"HTTPSession::XMLRequest response",
			"content",
			responseBody,
			"status code",
			httpResponse.StatusCode,
		)

		err = json.NewDecoder(bytes.NewReader(responseBody)).Decode(params.Response)
		if err != nil {
			return nil, fmt.Errorf("error decoding response: %v", err)
		}
	}
	return httpResponse, nil
}
