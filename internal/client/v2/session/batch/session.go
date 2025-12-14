package batch

import (
	"ksef/internal/certsdb"
	HTTP "ksef/internal/http"
)

type Session struct {
	httpClient *HTTP.Client
	certsDB    *certsdb.CertificatesDB
	workDir    string
}

func NewSession(
	httpClient *HTTP.Client,
	certsDB *certsdb.CertificatesDB,
	workDir string,
) *Session {
	return &Session{
		httpClient: httpClient,
		certsDB:    certsDB,
	}
}
