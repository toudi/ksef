package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"ksef/internal/interfaces"
	"ksef/internal/logging"
	"ksef/internal/runtime"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
	ErrUnableToCopyResponse = errors.New("unable to copy HTTP response to buffer")
	ErrUnexpectedBody       = errors.New("unexpected body (content type not specified and body is not a reader)")
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
	Vip                    *viper.Viper
}

func NewClient(host string) *Client {
	return &Client{Base: host}
}

func (rb *Client) Request(
	ctx context.Context,
	config RequestConfig,
	endpoint string,
) (*http.Response, error) {
	maxRetries := runtime.DefaultHttpRetries
	if rb.Vip != nil {
		maxRetries = runtime.HttpMaxRetries(rb.Vip)
	}

	var err error
	var resp *http.Response

	for numRetries := range maxRetries {
		resp, err = rb.requestAttempt(ctx, config, endpoint, numRetries, maxRetries)
		if err != nil {
			// let's check if the returned error code is telling us to slow down
			if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
				secondsToWait, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
				if secondsToWait > 0 {
					time.Sleep(time.Duration(secondsToWait) * time.Second)
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

func (rb *Client) SetRateLimiter(limiter *RequestRateLimit) {
	rb.rateLimiter = limiter
	if rb.Vip != nil {
		if err := rb.restoreRateLimitsState(rb.Vip); err != nil {
			logging.HTTPLogger.Error("Unable to restore rate limits state", "err", err)
		}
	}
}
