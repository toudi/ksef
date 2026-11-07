package ksef

import (
	"time"

	"ksef/internal/utils"
)

type InvoiceSubjectMetadata struct {
	NIP  string `json:"nip"`
	Name string `json:"name"`
}

// InvoiceMetadata represents the invoice metadata returned by the KSeF API.
// It is pure, implementation-agnostic data and is shared by the client,
// storage, and download layers.
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
	Metadata          map[string]any
}

func (im InvoiceMetadata) Checksum() string {
	return utils.Base64ToHex(im.InvoiceHashBase64)
}

func (im InvoiceMetadata) IssueTime() time.Time {
	issueTime, _ := time.Parse(time.DateOnly, im.IssueDate)
	return issueTime
}
