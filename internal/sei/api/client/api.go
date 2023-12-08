package client

import (
	"errors"
	"fmt"
	"ksef/internal/encryption"
)

const (
	ProductionEnvironment string = "prod"
	TestEnvironment       string = "test"
)

type environmentConfigType struct {
	Host         string
	RsaPublicKey string
}

var environmentConfig = map[string]environmentConfigType{
	ProductionEnvironment: {Host: "ksef.mf.gov.pl", RsaPublicKey: "klucze/prod/publicKey.pem"},
	TestEnvironment:       {Host: "ksef-test.mf.gov.pl", RsaPublicKey: "klucze/test/publicKey.pem"},
}

type EncryptionType struct {
	Cipher             *encryption.Cipher
	CipherTemplateVars encryption.CipherTemplateVarsType
}

var errUnknownEnvironment = errors.New("unknown environment")

type APIClient struct {
	Environment      environmentConfigType
	EnvironmentAlias string
	encryption       *EncryptionType
}

func APIClient_Init(environment string) (*APIClient, error) {
	if config, exists := environmentConfig[environment]; exists {
		return &APIClient{Environment: config, EnvironmentAlias: environment}, nil
	}

	return nil, errUnknownEnvironment
}

func (a *APIClient) Encryption() (*EncryptionType, error) {
	var err error

	if a.encryption == nil {
		a.encryption = &EncryptionType{}
		if a.encryption.Cipher, err = encryption.CipherInit(32); err != nil {
			return nil, fmt.Errorf("unable to initialize cipher: %v", err)
		}

		a.encryption.CipherTemplateVars.IV = make([]byte, len(a.encryption.Cipher.IV))
		copy(a.encryption.CipherTemplateVars.IV, a.encryption.Cipher.IV)
		a.encryption.CipherTemplateVars.EncryptionKey = make([]byte, len(a.encryption.Cipher.Key))
		copy(a.encryption.CipherTemplateVars.EncryptionKey, a.encryption.Cipher.Key)
	}

	return a.encryption, nil
}
