package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"ksef/core/ratelimits"
	"ksef/internal/interfaces"
	"ksef/internal/logging"
	"ksef/internal/runtime"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
	ErrUnableToCopyResponse = errors.New("unable to copy HTTP response to buffer")
	ErrUnexpectedBody       = errors.New("unexpected body (content type not specified and body is not a reader)")
	ErrUnwindingBody        = errors.New("unable to rewind body reader")
)

const (
	JSON = "application/json"
	XML  = "application/xml"
	BIN  = "application/octet-stream"
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
	OperationId     string
}

type Client struct {
	Base                   string
	rateLimiter            *RequestRateLimit
	AuthTokenRetrieverFunc interfaces.TokenRetrieverFunc
	RateLimitsDiscoverFunc interfaces.RateLimitsDiscoverFunc

	// RateLimitsState persists the rate limiter state between runs. When nil,
	// the state is neither restored nor saved.
	RateLimitsState ratelimits.StateStorage
	// MaxRetries is the number of attempts a request is retried. Values <= 0
	// fall back to runtime.DefaultHttpRetries.
	MaxRetries int
}

func NewClient(host string, opts ...Option) *Client {
	client := &Client{Base: host}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func (rb *Client) Request(
	ctx context.Context,
	config RequestConfig,
	endpoint string,
) (*http.Response, error) {
	maxRetries := rb.MaxRetries
	if maxRetries <= 0 {
		maxRetries = runtime.DefaultHttpRetries
	}

	var err error
	var resp *http.Response

	for numRetries := range maxRetries {
		// if the request is a bytes.Buffer then we need to temporarily
		// read it to memory so that we can rewind it after failed request.
		// otherwise the content would be lost since bytes.Buffer implements
		// Reader but not Seeker
		if body, isBuffer := config.Body.(*bytes.Buffer); isBuffer {
			logging.HTTPLogger.Debug("replace body with reader in order to be able to rewind it in case the request fails")
			config.Body = bytes.NewReader(body.Bytes())
		}
		resp, err = rb.requestAttempt(ctx, config, endpoint, numRetries, maxRetries)
		if err != nil {
			// let's check if the returned error code is telling us to slow down
			if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
				secondsToWait, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
				if secondsToWait > 0 {
					time.Sleep(time.Duration(secondsToWait) * time.Second)
				}
			}
			if config.Body != nil {
				logging.HTTPLogger.Debug("request has failed but it has a body. trying to rewind it so it can be resent")
				if seeker, ok := config.Body.(io.Seeker); ok {
					_, err := seeker.Seek(0, io.SeekStart)
					if err != nil {
						return nil, errors.Join(ErrUnwindingBody, err)
					}
				}
			}
			continue
		}
		// if we're here it means that the error was nil
		if rb.rateLimiter != nil {
			rb.rateLimiter.replaceLastEntry(config.OperationId)
		}
		break
	}

	rb.tryToPersistRateLimiterState(ctx)

	return resp, err
}

func (rb *Client) Download(ctx context.Context, url string, dest io.Writer) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(dest, resp.Body)
	return err
}

func jsonBodyReader(body any) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(body)
	return &buffer, err
}
