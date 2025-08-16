package v2

import (
	"context"
	"errors"
	"ksef/internal/sei/api/environment"
	"net/http"
)

const apiPrefix = "/api/v2"

var (
	errUnknownEnvironment = errors.New("unknown environment")
)

type APIClient struct {
	httpClient        *http.Client
	environmentConfig environment.EnvironmentConfig
	ctx               context.Context
}

func NewClient(ctx context.Context, environmentType environment.EnvironmentType) (*APIClient, error) {
	config, exists := environment.Environments[environmentType]
	if !exists {
		return nil, errUnknownEnvironment
	}
	return &APIClient{
		httpClient:        http.DefaultClient,
		environmentConfig: config,
		ctx:               ctx,
	}, nil
}
