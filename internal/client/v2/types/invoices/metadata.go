package invoices

import (
	"ksef/internal/utils"
	"time"
)

type InvoiceSubjectMetadata struct {
	NIP  string `json:"nip"`
	Name string `json:"name"`
}
type InvoiceMetadata struct {
	KSeFNumber        string                 `json:"ksefNumber"`
	InvoiceNumber     string                 `json:"invoiceNumber"`
	InvoiceType       string                 `json:"invoiceType"`
	IssueDate         string                 `json:"issueDate"`
	StorageDate       time.Time              `json:"permanentStorageDate"`
	Seller            InvoiceSubjectMetadata `json:"seller"`
	Buyer             InvoiceSubjectMetadata `json:"buyer"`
	InvoiceHashBase64 string                 `json:"invoiceHash"`
	Offline           bool
	Metadata          map[string]string
}

func (im InvoiceMetadata) Checksum() string {
	return utils.Base64ToHex(im.InvoiceHashBase64)
}

func (im InvoiceMetadata) IssueTime() time.Time {
	issueTime, _ := time.Parse(time.DateOnly, im.IssueDate)
	return issueTime
}

type InvoiceMetadataResponse struct {
	HasMore  bool              `json:"hasMore"`
	Invoices []InvoiceMetadata `json:"invoices"`
}
