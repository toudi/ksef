package invoices

import "ksef/core/ksef"

type InvoiceMetadataResponse struct {
	HasMore  bool                   `json:"hasMore"`
	Invoices []ksef.InvoiceMetadata `json:"invoices"`
}
