package resources

import (
	"ksef/core/certsdb"
	"ksef/core/certsdb/fs"
	"ksef/internal/runtime"

	"github.com/spf13/viper"
)

var (
	fsDataStorageInstance *fs.FSDataStorage
	fsMetaStorageInstance *fs.FSCertsStorage
	certsDBInstance       *certsdb.CertsDB
)

// CertsDataStorage used to be a private (unexported) struct, however
// because we also have a command that talks to the implementation
// directly (i.e. the filesystem-one) it needs to be exported.
func CertsDataStorage(vip *viper.Viper) (*fs.FSDataStorage, error) {
	if fsDataStorageInstance != nil {
		return fsDataStorageInstance, nil
	}

	keyring, err := Keyring(vip)
	if err != nil {
		return nil, err
	}

	data, err := fs.NewFSDataStorage(runtime.GetEnvironmentId(vip), keyring)
	if err != nil {
		return nil, err
	}
	fsDataStorageInstance = data

	// since these are simply small utility functions that talk to
	// the filesystem, there's no destructor required.

	return data, nil
}

// CertsMetaStorage is the filesystem-based implementation of certificates
// metadata storage.
func CertsMetaStorage(vip *viper.Viper) (*fs.FSCertsStorage, error) {
	if fsMetaStorageInstance != nil {
		return fsMetaStorageInstance, nil
	}

	data, err := CertsDataStorage(vip)
	if err != nil {
		return nil, err
	}

	meta, err := fs.NewFSCertStorage(runtime.GetEnvironmentId(vip), data)
	if err != nil {
		return nil, err
	}
	fsMetaStorageInstance = meta
	destructors = append(destructors, meta.Close)

	return meta, nil
}

// CertsDB returns the process-wide certificate database, wiring the metadata
// and data storages and applying the runtime-provided configuration.
func CertsDB(vip *viper.Viper) (*certsdb.CertsDB, error) {
	if certsDBInstance != nil {
		return certsDBInstance, nil
	}

	metaStorage, err := CertsMetaStorage(vip)
	if err != nil {
		return nil, err
	}
	dataStorage, err := CertsDataStorage(vip)
	if err != nil {
		return nil, err
	}

	db, err := certsdb.Init(
		metaStorage,
		dataStorage,
		certsdb.WithPreferredCertProfile(runtime.GetPreferredCertProfile(vip)),
	)
	if err != nil {
		return nil, err
	}
	certsDBInstance = db
	destructors = append(destructors, db.Close)

	return db, nil
}
