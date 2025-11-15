package certificates

import (
	"ksef/internal/certsdb"
	"ksef/internal/http"
	"ksef/internal/runtime"
)

type Manager struct {
	httpClient *http.Client
	certsDB    *certsdb.CertificatesDB
	env        runtime.Gateway
}

func NewManager(httpClient *http.Client, certsDB *certsdb.CertificatesDB, env runtime.Gateway) *Manager {
	return &Manager{
		httpClient: httpClient,
		certsDB:    certsDB,
		env:        env,
	}
}

func (m *Manager) SaveDB() error {
	return m.certsDB.Save()
}
