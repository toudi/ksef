package certsdb

import (
	"context"
	"errors"
)

type CertsDB struct {
	metadataStorage  CertificateMetadataStorage
	dataStorage      CertificateDataStorage
	preferredProfile string
	// env is the KSeF gateway id (e.g. "demo") used to narrow certificate
	// lookups. It is required and validated in Init.
	env string
}

type CertsDBInitializer func(*CertsDB)

// WithPreferredCertProfile sets the preferred certificate profile.
// Basically, if there are multiple certificates assigned to a single NIP / entity
// (e.g. multiple employees within the same company) this allows to narrow down
// authentication / offline certificate to a single one
func WithPreferredCertProfile(name string) CertsDBInitializer {
	return func(c *CertsDB) {
		c.preferredProfile = name
	}
}

func Init(
	environment string,
	metadataStorage CertificateMetadataStorage,
	dataStorage CertificateDataStorage,
	initializers ...CertsDBInitializer,
) (*CertsDB, error) {
	if environment == "" {
		return nil, ErrEnvironmentRequired
	}

	certsDB := &CertsDB{
		env:             environment,
		metadataStorage: metadataStorage,
		dataStorage:     dataStorage,
	}
	for _, init := range initializers {
		init(certsDB)
	}

	return certsDB, nil
}

func (c *CertsDB) Close(ctx context.Context) error {
	return errors.Join(
		c.metadataStorage.Close(ctx),
		c.dataStorage.Close(ctx),
	)
}
