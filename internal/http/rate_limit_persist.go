package http

import (
	"ksef/internal/invoicesdb/config"
	"ksef/internal/logging"
	"ksef/internal/runtime"
	"ksef/internal/utils"
	ratelimiter "ksef/internal/utils/rate-limiter"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

const rateLimitsState = "rate-limits-state.yaml"

func (rb *Client) SaveRateLimitsState(vip *viper.Viper) error {
	logging.HTTPLogger.Debug("persisting rate limits state")
	if rb.rateLimiter == nil {
		return nil
	}

	state := make(map[string][]ratelimiter.LimiterSlotEntries)

	now := time.Now()

	for operationId, limiter := range rb.rateLimiter.limits {
		requestsWithinSlot := limiter.EntriesWithinSlot(now)
		if len(requestsWithinSlot) > 0 {
			state[operationId] = requestsWithinSlot
		}
	}

	outputFilename, err := rb.stateFilename(vip)
	if err != nil {
		return err
	}

	logging.HTTPLogger.Debug("rate limits state file", "filename", outputFilename)

	return utils.SaveYAML(state, outputFilename)
}

func (rb *Client) tryToPersistRateLimiterState() {
	if rb.Vip != nil {
		rb.SaveRateLimitsState(rb.Vip)
	}
}

func (rb *Client) restoreRateLimitsState(vip *viper.Viper) error {
	stateFilename, err := rb.stateFilename(vip)
	if err != nil {
		return err
	}

	state := make(map[string][]ratelimiter.LimiterSlotEntries)
	fileReader, exists, _ := utils.FileExists(stateFilename)
	if !exists {
		return nil
	}
	defer fileReader.Close()
	if err = utils.ReadYAML(fileReader, &state); err != nil {
		return err
	}

	for operationId, requestsWithinSlot := range state {
		limiter, exists := rb.rateLimiter.limits[operationId]
		if !exists {
			logging.HTTPLogger.Warn("rate limiter for operation ID does not exist", "operationId", operationId)
			continue
		}
		limiter.LoadEntries(requestsWithinSlot)
	}

	return nil
}

func (rb *Client) stateFilename(vip *viper.Viper) (string, error) {
	cfg := config.GetInvoicesDBConfig(vip)

	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return "", err
	}

	prefix := filepath.Join(
		cfg.Root,
		runtime.GetEnvironmentId(vip),
		nip,
	)

	outputFilename := filepath.Join(
		prefix,
		rateLimitsState,
	)

	return outputFilename, nil
}
