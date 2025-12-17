package status

import "ksef/internal/client/v2/upo"

const (
	SessionStatusProcessed int = 200
)

type StatusResponse struct {
	Status struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
	} `json:"status" yaml:"status,omitempty"`
	Upo struct {
		Pages []upo.UPODownloadPage `json:"pages"`
	} `json:"upo" yaml:"-"`
	InvoiceCount           int `json:"invoiceCount" yaml:"invoice-count"`
	SuccessfulInvoiceCount int `json:"successfulInvoiceCount" yaml:"successful-invoice-count"`
	FailedInvoiceCount     int `json:"failedInvoiceCount" yaml:"failed-invoice-count"`
}

func (s *StatusResponse) IsProcessed() bool {
	return s.Status.Code >= 200
}
