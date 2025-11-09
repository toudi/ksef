package token

import "ksef/internal/certsdb"

type initializerFunc = func(handler *TokenHandler)

func WithCertsDB(certsDB *certsdb.CertificatesDB) initializerFunc {
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
