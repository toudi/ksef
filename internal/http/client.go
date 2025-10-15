package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"ksef/internal/interfaces"
	"ksef/internal/logging"
	"net/http"
	"net/url"
	"time"
)

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
	ErrUnableToCopyResponse = errors.New("unable to copy HTTP response to buffer")
	ErrUnexpectedBody       = errors.New("unexpected body (content type not specified and body is not a reader)")
)

const (
	JSON = "application/json"
	XML  = "application/xml"
)

type RequestConfig struct {
	Timeout         time.Duration
	Headers         map[string]string
	QueryParams     map[string]string
	ContentType     string
	Body            any
	Dest            any
	DestWriter      io.Writer
	DestContentType string
	ExpectedStatus  int
	Method          string
}

type Client struct {
	Base                   string
	AuthTokenRetrieverFunc interfaces.TokenRetrieverFunc
}

func (rb *Client) Request(ctx context.Context, config RequestConfig, endpoint string) (*http.Response, error) {
	var cancel context.CancelFunc

	if config.Timeout.Abs() == 0 {
		config.Timeout = 15 * time.Second
	}

	if config.ContentType == "" {
		config.ContentType = JSON
	}

	ctx, cancel = context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	fullUrl, err := url.Parse(rb.Base + "/" + endpoint)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	var logger = logging.HTTPLogger.With("method", config.Method, "url", fullUrl.String())

	if config.Body != nil {
		var isReader bool

		if config.ContentType == JSON {
			logger = logger.With("body", fmt.Sprintf("%+v", config.Body))
			body, err = jsonBodyReader(config.Body)
			if err != nil {
				return nil, err
			}
		} else {
			body, isReader = config.Body.(io.Reader)
			if !isReader {
				return nil, ErrUnexpectedBody
			}
		}
	}

	logger.Debug("request")

	req, err := http.NewRequestWithContext(ctx, config.Method, fullUrl.String(), body)
	if err != nil {
		return nil, err
	}
	if config.QueryParams != nil {
		params := req.URL.Query()
		for paramName, paramValue := range config.QueryParams {
			params.Set(paramName, paramValue)
		}
		req.URL.RawQuery = params.Encode()
	}

	if config.ContentType != "" {
		req.Header.Set("Content-Type", config.ContentType)
	}

	for headerName, headerValue := range config.Headers {
		req.Header.Set(headerName, headerValue)
	}

	if rb.AuthTokenRetrieverFunc != nil {
		token, err := rb.AuthTokenRetrieverFunc()
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return resp, err
	}

	var bodyBuffer bytes.Buffer
	if _, err = io.Copy(&bodyBuffer, resp.Body); err != nil {
		return nil, ErrUnableToCopyResponse
	}
	defer resp.Body.Close()

	logging.HTTPLogger.Debug(
		"response",
		"body",
		bodyBuffer.String(),
	)

	if config.ExpectedStatus > 0 && resp.StatusCode != config.ExpectedStatus {
		return resp, fmt.Errorf("%w: %d vs %d", ErrUnexpectedStatusCode, resp.StatusCode, config.ExpectedStatus)
	}

	if config.DestContentType == "" {
		// if no content type is specified, simply copy to dest
		if config.DestWriter != nil {
			_, err = io.Copy(config.DestWriter, &bodyBuffer)
			return resp, err
		}

		return resp, nil
	}

	if config.DestContentType == JSON {
		decoder := json.NewDecoder(&bodyBuffer)
		return resp, decoder.Decode(config.Dest)
	}

	return resp, nil
}

func jsonBodyReader(body any) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	var encoder = json.NewEncoder(&buffer)
	err := encoder.Encode(body)
	return &buffer, err
}
