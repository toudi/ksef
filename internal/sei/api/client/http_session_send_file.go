package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

type SendFileParams struct {
	Method       string
	Endpoint     string
	FileName     string
	JSONResponse interface{}
	Logger       *slog.Logger
}

func (hs *HTTPSession) SendFile(params SendFileParams) (*http.Response, error) {
	var log *slog.Logger = params.Logger

	fileReader, err := os.Open(params.FileName)
	if err != nil {
		return nil, fmt.Errorf("unable to open file for reading: %v", err)
	}
	defer fileReader.Close()

	request, err := hs.Request(params.Method, params.Endpoint, fileReader, log)
	log.Debug("HTTPSEssion::SendFile", "method", params.Method, "endpoint", params.Endpoint)

	if err != nil {
		return nil, fmt.Errorf("error preparing request: %v", err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error sending file: %v", err)
	}
	defer response.Body.Close()

	log.Debug(
		"HTTPSession::SendFile response status",
		"method", params.Method,
		"endpoint", params.Endpoint,
		"status code", response.StatusCode,
	)

	if response.StatusCode/100 != 2 {
		responseContent, err := io.ReadAll(response.Body)
		log.Debug(
			"HTTPSession::SendFile response body",
			"method", params.Method,
			"endpoint", params.Endpoint,
			"content", string(responseContent),
		)
		if err != nil {
			return nil, fmt.Errorf("error reading response:%v", err)
		}
		return nil, fmt.Errorf(
			"unexpected response code from initResponse:\n%s\n",
			string(responseContent),
		)
	}

	err = json.NewDecoder(response.Body).Decode(params.JSONResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %v", err)
	}

	return response, nil
}
