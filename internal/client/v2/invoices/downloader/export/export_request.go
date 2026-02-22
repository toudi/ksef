package export

import (
	downloadertypes "ksef/internal/client/v2/invoices/downloader/types"
	"ksef/internal/encryption"
	"strconv"
	"strings"
)

const (
	endpointInvoicesExport       = "/v2/invoices/exports"
	endpointInvoicesExportStatus = "/v2/invoices/exports/%s"
)

const (
	exportStatusReady   = 200
	exportStatusExpired = 210
)

type exportRequest struct {
	Encryption encryption.CipherHTTPRequest       `json:"encryption"`
	Filters    downloadertypes.InvoiceListRequest `json:"filters"`
}

type exportResponse struct {
	ReferenceNumber string `json:"referenceNumber"`
}

type exportStatusPart struct {
	OrdinalNumber     int    `json:"ordinalNumber"`
	PartName          string `json:"partName"`
	HTTPMethod        string `json:"method"`
	URL               string `json:"url"`
	PartSize          int64  `json:"partSize"`
	PartHash          string `json:"partHash"`
	EncryptedPartSize int64  `json:"encryptedPartSize"`
	EncryptedPartHash string `json:"encryptedPartHash"`
}

func (p exportStatusPart) decryptedFilename() string {
	return strings.Replace(p.PartName, ".zip.aes", ".zip."+strconv.Itoa(p.OrdinalNumber), 1)
}

type exportStatusResponse struct {
	Status struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
	} `json:"status"`

	Package struct {
		InvoiceCount int64              `json:"invoiceCount"`
		Size         int64              `json:"size"`
		Parts        []exportStatusPart `json:"parts"`
	} `json:"package"`

	IsTruncated bool `json:"isTruncated"`
}

func (sr *exportStatusResponse) Ready() bool {
	return sr.Status.Code == exportStatusReady
}

func (sr *exportStatusResponse) Invalid() (bool, string) {
	return sr.Status.Code/100 > 2, sr.Status.Description
}
