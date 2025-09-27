package invoices

type InvoiceSellerMetadata struct {
	Identifier string `json:"identifier"`
}
type InvoiceMetadata struct {
	KSeFNumber    string                `json:"ksefNumber"`
	InvoiceNumber string                `json:"invoiceNumber"`
	InvoiceType   string                `json:"invoiceType"`
	Seller        InvoiceSellerMetadata `json:"seller"`
}

type InvoiceMetadataResponse struct {
	HasMore  bool              `json:"hasMore"`
	Invoices []InvoiceMetadata `json:"invoices"`
}
