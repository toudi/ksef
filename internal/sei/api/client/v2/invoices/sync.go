package invoices

import (
	"context"
	"ksef/internal/http"
	"ksef/internal/registry"
	baseHttp "net/http"
	"strconv"
	"time"
)

type DateRangeType string

const (
	endpointInvoicesMetadata               = "/api/v2/invoices/query/metadata"
	DateRangeTypeIssue       DateRangeType = "Issue"
)

type InvoiceMetadata struct {
	KSeFNumber    string `json:"ksefNumber"`
	InvoiceNumber string `json:"invoiceNumber"`
}

type InvoiceMetadataResponse struct {
	HasMore  bool              `json:"hasMore"`
	Invoices []InvoiceMetadata `json:"invoices"`
}

type InvoiceMetadataRequest struct {
	SubjectType SubjectType `json:"subjectType"`
	DateRange   struct {
		DateType DateRangeType `json:"dateType"`
		From     time.Time     `json:"from"`
	} `json:"dateRange"`
}

func Sync(ctx context.Context, httpClient *http.Client, params SyncParams, registry *registry.InvoiceRegistry) error {
	var (
		finished bool
		page     int
		req      InvoiceMetadataRequest
		resp     InvoiceMetadataResponse
		err      error
	)

	req.DateRange.DateType = DateRangeTypeIssue
	req.DateRange.From = registry.QueryCriteria.DateFrom

	for !finished {
		_, err = httpClient.Request(
			ctx,
			http.RequestConfig{
				Method: baseHttp.MethodPost,
				QueryParams: map[string]string{
					"pageOffset": strconv.Itoa(page),
					"pageSize":   strconv.Itoa(params.PageSize),
				},
				Body:            req,
				ContentType:     http.JSON,
				Dest:            &resp,
				DestContentType: http.JSON,
				ExpectedStatus:  baseHttp.StatusOK,
			},
			endpointInvoicesMetadata,
		)
		if err != nil {
			return err
		}

		for _, invoice := range resp.Invoices {
			if registry.Contains(invoice.KSeFNumber) {
				continue
			}
			// TODO: download and print PDF
		}

		page += 1
	}

	return registry.Save("")
}
