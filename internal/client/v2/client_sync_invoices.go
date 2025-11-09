package v2

import (
	"context"
	"ksef/internal/client/v2/invoices"
)

func (c *APIClient) SyncInvoices(ctx context.Context, params invoices.SyncParams) error {
	return invoices.Sync(
		ctx,
		c.authenticatedHTTPClient(),
		params,
		c.registry,
	)
}
