package api

import (
	"errors"
	"fmt"
	"ksef/common/aes"
)

const (
	ProductionEnvironment string = "prod"
	TestEnvironment       string = "test"
	StatusFileFormatYAML  string = "yaml"
	StatusFileFormatJSON  string = "json"
)

type environmentConfigType struct {
	host         string
	rsaPublicKey string
}

type cipherTemplateVarsType struct {
	IV            []byte
	EncryptionKey []byte
}

var environmentConfig = map[string]environmentConfigType{
	ProductionEnvironment: {host: "ksef.mf.gov.pl", rsaPublicKey: "klucze/prod/publicKey.pem"},
	TestEnvironment:       {host: "ksef-test.mf.gov.pl", rsaPublicKey: "klucze/test/publicKey.pem"},
}

var errUnknownEnvironment = errors.New("unknown environment")

type API struct {
	environment        environmentConfigType
	environmentAlias   string
	cipher             *aes.Cipher
	cipherTemplateVars cipherTemplateVarsType
	requestFactory     *RequestFactory
}

func API_Init(environment string) (*API, error) {
	if config, exists := environmentConfig[environment]; exists {
		var err error

		api := &API{environment: config, environmentAlias: environment}

		if api.cipher, err = aes.CipherInit(32); err != nil {
			return nil, fmt.Errorf("unable to initialize cipher: %v", err)
		}

		api.cipherTemplateVars.IV = make([]byte, len(api.cipher.IV))
		copy(api.cipherTemplateVars.IV, api.cipher.IV)
		api.cipherTemplateVars.EncryptionKey = make([]byte, len(api.cipher.Key))
		copy(api.cipherTemplateVars.EncryptionKey, api.cipher.Key)

		api.requestFactory = NewRequestFactory(api)

		return api, nil
	}

	return nil, errUnknownEnvironment
}
