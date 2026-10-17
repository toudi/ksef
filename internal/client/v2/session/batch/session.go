package batch

import (
	"ksef/core/certsdb"
	HTTP "ksef/internal/http"
)

type Session struct {
	httpClient *HTTP.Client
	certsDB    *certsdb.CertsDB
	workDir    string
}

func NewSession(
	httpClient *HTTP.Client,
	certsDB *certsdb.CertsDB,
	workDir string,
) *Session {
	return &Session{
		httpClient: httpClient,
		certsDB:    certsDB,
	}
}
