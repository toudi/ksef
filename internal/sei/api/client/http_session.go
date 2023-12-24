package client

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type HTTPSession struct {
	host    string
	headers http.Header
}

func NewHTTPSession(host string) *HTTPSession {
	return &HTTPSession{host: host, headers: make(http.Header)}
}

func (hs *HTTPSession) SetHeader(header string, value string) {
	hs.headers.Add(header, value)
}

func (hs *HTTPSession) Request(
	method string,
	endpoint string,
	payload io.Reader,
	log *slog.Logger,
) (*http.Request, error) {
	url, err := url.Parse(fmt.Sprintf("https://%s%s", hs.host, endpoint))
	if err != nil {
		return nil, fmt.Errorf("unable to parse URL: %v", err)
	}
	request, err := http.NewRequest(method, url.String(), payload)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	log.Debug("HTTPSession::Request", "method", method, "URL", url.String())

	for key, values := range hs.headers {
		log.Debug("HTTPSession::Request add header", "header", key, "value", values[0])
		request.Header.Add(key, values[0])
	}

	return request, nil
}
