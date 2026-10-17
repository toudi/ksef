package certificates

import (
	"context"
	"ksef/core/certsdb"
	"ksef/internal/http"
)

type Manager struct {
	httpClient *http.Client
	certsDB    *certsdb.CertsDB
	env        string
}

func NewManager(httpClient *http.Client, certsDB *certsdb.CertsDB, env string) *Manager {
	return &Manager{
		httpClient: httpClient,
		certsDB:    certsDB,
		env:        env,
	}
}

func (m *Manager) SaveDB(ctx context.Context) error {
	return m.certsDB.Close(ctx)
}
