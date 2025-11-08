package invoices

type InvoiceSubjectMetadata struct {
	NIP  string `json:"nip"`
	Name string `json:"name"`
}
type InvoiceMetadata struct {
	KSeFNumber    string                 `json:"ksefNumber"`
	InvoiceNumber string                 `json:"invoiceNumber"`
	InvoiceType   string                 `json:"invoiceType"`
	IssueDate     string                 `json:"issueDate"`
	Seller        InvoiceSubjectMetadata `json:"seller"`
	Buyer         InvoiceSubjectMetadata `json:"buyer"`
	InvoiceHash   string                 `json:"invoiceHash"`
	Offline       bool
}

type InvoiceMetadataResponse struct {
	HasMore  bool              `json:"hasMore"`
	Invoices []InvoiceMetadata `json:"invoices"`
}
