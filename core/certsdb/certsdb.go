package certsdb

import "context"

type certsDB struct {
	storage CertificateStorage
}

func Init(storage CertificateStorage) *certsDB {
	return &certsDB{
		storage: storage,
	}
}

func (c *certsDB) Close(ctx context.Context) error {
	return c.storage.Close(ctx)
}
