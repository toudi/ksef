package invoices

type InvoiceSellerMetadata struct {
	NIP string `json:"nip"`
}
type InvoiceMetadata struct {
	KSeFNumber    string                `json:"ksefNumber"`
	InvoiceNumber string                `json:"invoiceNumber"`
	InvoiceType   string                `json:"invoiceType"`
	IssueDate     string                `json:"issueDate"`
	Seller        InvoiceSellerMetadata `json:"seller"`
	Offline       bool
}

type InvoiceMetadataResponse struct {
	HasMore  bool              `json:"hasMore"`
	Invoices []InvoiceMetadata `json:"invoices"`
}
