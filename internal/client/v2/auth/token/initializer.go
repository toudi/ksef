package token

import (
	"ksef/core/certsdb"
	kr "ksef/core/keyring"
)

type initializerFunc = func(handler *TokenHandler)

func WithCertsDB(certsDB *certsdb.CertsDB) initializerFunc {
	return func(handler *TokenHandler) {
		handler.certsDB = certsDB
	}
}

func WithDumpChallenge(output string) initializerFunc {
	return func(handler *TokenHandler) {
		handler.mode = modeDumpChallenge
		handler.challengeDumpPath = output
	}
}

func WithSignedChallengeFile(signedChallengeFile string) initializerFunc {
	return func(handler *TokenHandler) {
		handler.mode = modeUseSignedFile
		handler.signedChallengeFile = signedChallengeFile
	}
}

func WithKeyring(keyring kr.Keyring) initializerFunc {
	return func(handler *TokenHandler) {
		handler.keyring = keyring
	}
}
