package ratelimits

import (
	"context"
	"ksef/internal/http"
	"ksef/internal/logging"
	ratelimiter "ksef/internal/utils/rate-limiter"
	baseHTTP "net/http"
	"time"
)

const endpointRateLimits = "/api/v2/rate-limits"

type RateLimitsEntry struct {
	PerSecond int `json:"perSecond"`
	PerMinute int `json:"perMinute"`
	PerHour   int `json:"perHour"`
}

type RateLimitsResponse map[string]RateLimitsEntry

func (r *RateLimitsResponse) ToRequestRateLimiter() map[string]*ratelimiter.Limiter {
	limits := make(map[string]*ratelimiter.Limiter)
	for requestId, limitsDef := range *r {
		limits[requestId] = ratelimiter.NewLimiter(
			[]ratelimiter.RateLimit{
				{
					Slot:  time.Second,
					Limit: limitsDef.PerSecond,
				},
				{
					Slot:  time.Minute,
					Limit: limitsDef.PerMinute,
				},
				{
					Slot:  time.Hour,
					Limit: limitsDef.PerHour,
				},
			},
		)
	}
	return limits
}

func DiscoverRateLimits(
	ctx context.Context,
	host string,
	authToken string,
) (map[string]*ratelimiter.Limiter, error) {
	var resp RateLimitsResponse
	tmpClient := http.NewClient(host)
	_, err := tmpClient.Request(
		ctx,
		http.RequestConfig{
			Headers:         map[string]string{"Authorization": "Bearer " + authToken},
			Dest:            &resp,
			DestContentType: http.JSON,
			Method:          baseHTTP.MethodGet,
			ExpectedStatus:  baseHTTP.StatusOK,
		},
		endpointRateLimits,
	)
	if err != nil {
		return nil, err
	}
	// we got a response. now we need to translate this to our internal rate limits structure.
	logging.HTTPLogger.Debug("discovered rate limits", "limits", resp)
	requestRateLimiter := resp.ToRequestRateLimiter()
	return requestRateLimiter, nil
}
