package v2

import (
	"ksef/core/certsdb"
	kr "ksef/core/keyring"
	"ksef/core/ratelimits"
	"ksef/internal/client/v2/auth/validator"
)

type InitializerFunc func(c *APIClient)

func WithAuthValidator(validator validator.AuthChallengeValidator) func(client *APIClient) {
	return func(client *APIClient) {
		client.authChallengeValidator = validator
	}
}

func WithCertificatesDB(certsDB *certsdb.CertsDB) func(client *APIClient) {
	return func(client *APIClient) {
		client.certsDB = certsDB
	}
}

func WithKeyring(keyring kr.Keyring) func(client *APIClient) {
	return func(client *APIClient) {
		client.keyring = keyring
	}
}

func WithRateLimitsStateStorage(storage ratelimits.StateStorage) func(client *APIClient) {
	return func(client *APIClient) {
		client.rateLimitsState = storage
	}
}

func WithoutTokenManager() func(client *APIClient) {
	return func(client *APIClient) {
		client.runTokenManager = false
	}
}
