package certificates

import (
	"ksef/internal/certsdb"
	"ksef/internal/config"
	"ksef/internal/http"
)

type Manager struct {
	httpClient *http.Client
	certsDB    *certsdb.CertificatesDB
	env        config.Gateway
}

func NewManager(httpClient *http.Client, certsDB *certsdb.CertificatesDB, env config.Gateway) *Manager {
	return &Manager{
		httpClient: httpClient,
		certsDB:    certsDB,
		env:        env,
	}
}

func (m *Manager) SaveDB() error {
	return m.certsDB.Save()
}
