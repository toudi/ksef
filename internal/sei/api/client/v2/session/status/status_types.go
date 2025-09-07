package status

import "ksef/internal/sei/api/client/v2/upo"

const (
	SessionStatusProcessed int = 200
)

type StatusResponse struct {
	Status struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
	} `json:"status"`
	Upo struct {
		Pages []upo.UPODownloadPage `json:"pages"`
	} `json:"upo"`
	InvoiceCount           int `json:"invoiceCount"`
	SuccessfulInvoiceCount int `json:"successfulInvoiceCount"`
	FailedInvoiceCount     int `json:"failedInvoiceCount"`
}
