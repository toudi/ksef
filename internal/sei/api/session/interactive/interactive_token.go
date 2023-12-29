package interactive

import "github.com/zalando/go-keyring"

func (i *InteractiveSession) retrieveGateweayToken(issuer string) (string, error) {
	return keyring.Get(i.apiClient.Environment.Host, issuer)
}

func (i *InteractiveSession) PersistToken(issuer string, token string) error {
	return keyring.Set(i.apiClient.Environment.Host, issuer, token)
}
