package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type JSONRequestParams struct {
	Method   string
	Endpoint string
	Payload  interface{}
	Response interface{}
	Logger   *slog.Logger
}

func (hs *HTTPSession) JSONRequest(params JSONRequestParams) (*http.Response, error) {
	var encodedPayload bytes.Buffer
	var log *slog.Logger = params.Logger
	if err := json.NewEncoder(&encodedPayload).Encode(params.Payload); err != nil {
		return nil, fmt.Errorf("error encoding JSON: %v", err)
	}

	request, err := hs.Request(params.Method, params.Endpoint, &encodedPayload, log)

	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")

	httpResponse, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error performing HTTP request: %v", err)
	}
	defer httpResponse.Body.Close()

	httpResponseBody, _ := io.ReadAll(httpResponse.Body)
	log.Debug(
		"HTTPSession::JSONRequest response",
		"content",
		string(httpResponseBody),
		"status code",
		httpResponse.StatusCode,
	)

	if params.Response != nil {
		err = json.NewDecoder(bytes.NewReader(httpResponseBody)).Decode(params.Response)
		if err != nil {
			return nil, fmt.Errorf("error decoding response: %v", err)
		}
	}
	return httpResponse, nil
}
