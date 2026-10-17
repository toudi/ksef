package certsdb

import "context"

type CertsDB struct {
	storage          CertificateStorage
	preferredProfile string
	env              string
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

// WithEnvironment sets the environment that certificate selection operates on.
func WithEnvironment(id string) CertsDBInitializer {
	return func(c *CertsDB) {
		c.env = id
	}
}

func Init(storage CertificateStorage, initializers ...CertsDBInitializer) (*CertsDB, error) {
	certsDB := &CertsDB{
		storage: storage,
	}
	for _, init := range initializers {
		init(certsDB)
	}
	return certsDB, nil
}

func (c *CertsDB) Close(ctx context.Context) error {
	return c.storage.Close(ctx)
}
