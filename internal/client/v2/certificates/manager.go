package certificates

import (
	"ksef/internal/certsdb"
	"ksef/internal/http"
)

type Manager struct {
	httpClient *http.Client
	certsDB    *certsdb.CertificatesDB
	env        string
}

func NewManager(httpClient *http.Client, certsDB *certsdb.CertificatesDB, env string) *Manager {
	return &Manager{
		httpClient: httpClient,
		certsDB:    certsDB,
		env:        env,
	}
}

func (m *Manager) SaveDB() error {
	return m.certsDB.Save()
}
