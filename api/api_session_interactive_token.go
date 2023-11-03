package api

import "github.com/zalando/go-keyring"

func (i *InteractiveSession) retrieveGateweayToken(issuer string) (string, error) {
	return keyring.Get(i.api.environment.host, issuer)
}

func (i *InteractiveSession) PersistToken(issuer string, token string) error {
	return keyring.Set(i.api.environment.host, issuer, token)
}
