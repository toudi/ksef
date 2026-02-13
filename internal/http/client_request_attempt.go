package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"ksef/internal/logging"
	"net/http"
	"net/url"
	"time"
)

func (rb *Client) requestAttempt(
	ctx context.Context,
	config RequestConfig,
	endpoint string,
	attempt int,
	maxAttempts int,
) (*http.Response, error) {
	fullUrl, err := url.Parse(rb.Base + endpoint)
	if err != nil {
		return nil, err
	}
	method := http.MethodGet
	if config.Method != "" {
		method = config.Method
	}
	logger := logging.HTTPLogger.With(
		"method", method, "url", fullUrl.String(), "attempt", attempt+1, "max", maxAttempts,
		"operationId", config.OperationId,
	)

	// retrieve auth token if required:
	var bearerToken string
	if rb.AuthTokenRetrieverFunc != nil {
		token, err := rb.AuthTokenRetrieverFunc()
		if err != nil {
			return nil, err
		}
		bearerToken = token

		if rb.rateLimiter == nil && rb.RateLimitsDiscoverFunc != nil {
			// try to discover rate limits
			rateLimits, err := rb.RateLimitsDiscoverFunc(ctx, rb.Base, token)
			if err != nil {
				logging.HTTPLogger.Error("Unable to discover rate limits", "err", err)
			} else {
				rb.SetRateLimiter(
					NewRequestRateLimit(logging.HTTPLogger, rateLimits),
				)
			}
		}
	}

	// call rate limiter if possible - but before initializing context with timeout.
	// otherwise the request timeout would hit due to rate limiting.
	if rb.rateLimiter != nil {
		logger.Debug("calling rateLimiter.Wait()", "operationId", config.OperationId)
		rb.rateLimiter.Wait(config.OperationId)
	}

	if config.Timeout.Abs() == 0 {
		config.Timeout = 15 * time.Second
	}

	if config.ContentType == "" {
		config.ContentType = JSON
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	var body io.Reader

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

	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		rb.tryToPersistRateLimiterState()
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
		if resp.StatusCode == http.StatusTooManyRequests {
			rb.tryToPersistRateLimiterState()
		}
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
