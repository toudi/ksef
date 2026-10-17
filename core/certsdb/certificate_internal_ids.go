package certsdb

import (
	"context"
	"iter"
	"strings"

	"github.com/samber/lo"
)

func (c Certificate) GetInternalIDForNIP(nip string) *string {
	for _, internalId := range c.InternalIDs {
		if strings.HasPrefix(internalId, nip) {
			return lo.ToPtr(internalId)
		}
	}

	return nil
}

func (c *CertsDB) AddInternalIDsToCertificate(
	ctx context.Context,
	certificateId string,
	internalIds iter.Seq[string],
) error {
	return c.storage.AddInternalIDsToCertificate(
		ctx,
		certificateId,
		internalIds,
	)
}

func (c *CertsDB) RemoveInternalIdsFromCertificate(
	ctx context.Context,
	certificateId string,
	internalIds iter.Seq[string],
) error {
	return c.storage.RemoveInternalIDsFromCertificate(
		ctx,
		certificateId,
		internalIds,
	)
}
