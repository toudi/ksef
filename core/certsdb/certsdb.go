package certsdb

import "context"

type CertsDB struct {
	storage CertificateStorage
}

func Init(storage CertificateStorage) *CertsDB {
	return &CertsDB{
		storage: storage,
	}
}

func (c *CertsDB) Close(ctx context.Context) error {
	return c.storage.Close(ctx)
}
