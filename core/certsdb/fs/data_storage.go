package fs

import (
	"context"
	"ksef/core/certsdb"
	"ksef/core/keyring"
	"os"
	"path/filepath"
)

// FSDataStorage is the implementation of .. well .. the certificates
// and private keys data storage (their bytes). it takes the keyring
// as dependency due to encryption
type FSDataStorage struct {
	envId string
	kr    keyring.Keyring
}

var _ certsdb.CertificateDataStorage = (*FSDataStorage)(nil)

// NewFSDataStorage creates the data storage over the certificates directory.
func NewFSDataStorage(environmentId string, kr keyring.Keyring) (*FSDataStorage, error) {
	if err := os.MkdirAll(certificatesDir, 0775); err != nil {
		return nil, err
	}
	return &FSDataStorage{envId: environmentId, kr: kr}, nil
}

// Close is a no-op - the keyring is owned by the caller and closed upstream.
func (d *FSDataStorage) Close(ctx context.Context) error {
	return nil
}

func (d *FSDataStorage) SetEnvironmentId(envId string) {
	d.envId = envId
}

// fileName resolves the file for the given id.
func (d *FSDataStorage) fileName(id string, suffix string) string {
	return filepath.Join(certificatesDir, d.envId+"-"+id+suffix)
}
