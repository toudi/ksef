package http

import (
	"context"
	"time"

	"ksef/internal/logging"
)

func (rb *Client) SetRateLimiter(ctx context.Context, limiter *RequestRateLimit) {
	rb.rateLimiter = limiter
	if rb.RateLimitsState == nil {
		return
	}
	if err := rb.restoreRateLimitsState(ctx); err != nil {
		logging.HTTPLogger.Error("Unable to restore rate limits state", "err", err)
	}
}

func (rb *Client) saveRateLimitsState(ctx context.Context) error {
	if rb.rateLimiter == nil || rb.RateLimitsState == nil {
		return nil
	}

	logging.HTTPLogger.Debug("handing rate limits state over to the storage")

	return rb.RateLimitsState.Set(ctx, rb.rateLimiter.snapshot(time.Now()))
}

func (rb *Client) tryToPersistRateLimiterState(ctx context.Context) {
	if err := rb.saveRateLimitsState(ctx); err != nil {
		logging.HTTPLogger.Error("unable to persist rate limits state", "err", err)
	}
}

func (rb *Client) restoreRateLimitsState(ctx context.Context) error {
	state, err := rb.RateLimitsState.Load(ctx)
	if err != nil {
		return err
	}

	rb.rateLimiter.load(state)

	return nil
}

func (rb *Client) GetLimit(operationId string, slot time.Duration) (int, error) {
	if rb.rateLimiter == nil {
		return 0, errLimitNotFound
	}
	return rb.rateLimiter.GetLimit(operationId, slot)
}
