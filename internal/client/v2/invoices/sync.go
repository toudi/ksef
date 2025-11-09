package invoices

import (
	"context"
	"errors"
	types "ksef/internal/client/v2/types/invoices"
	"ksef/internal/http"
	"ksef/internal/pdf"
	"ksef/internal/pdf/printer"
	"ksef/internal/registry"
	baseHttp "net/http"
	"strconv"
	"time"
)

var (
	ErrUnableToInitializePDFPrintingEngine = errors.New("nie udało się zainicjować silnika wydruku PDF")
)

type DateRangeType string

const (
	endpointInvoicesMetadata               = "/api/v2/invoices/query/metadata"
	DateRangeTypeIssue       DateRangeType = "Issue"
)

type InvoiceMetadataRequest struct {
	SubjectType types.SubjectType `json:"subjectType"`
	DateRange   struct {
		DateType DateRangeType `json:"dateType"`
		From     time.Time     `json:"from"`
	} `json:"dateRange"`
}

func Sync(ctx context.Context, httpClient *http.Client, params SyncParams, registry *registry.InvoiceRegistry) error {
	var (
		finished       bool
		page           int
		req            InvoiceMetadataRequest
		resp           types.InvoiceMetadataResponse
		err            error
		printingEngine printer.PDFPrinter
	)

	req.SubjectType = params.SubjectType
	req.DateRange.DateType = DateRangeType(registry.Sync.QueryCriteria.DateType)
	req.DateRange.From = registry.Sync.QueryCriteria.DateFrom

	if params.PDF {
		printingEngine, err = pdf.GetLocalPrintingEngine()
		if err != nil {
			return ErrUnableToInitializePDFPrintingEngine
		}
	}

	downloader := NewInvoiceDownloader(httpClient, registry)

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
			var filename string
			var checksum string

			if registry.Contains(invoice.KSeFNumber) {
				continue
			}

			// download invoice to a local file
			if filename, checksum, err = downloader.Download(ctx, invoice, params.SubjectType); err != nil {
				return err
			}
			if params.PDF {
				registryMeta, err := registry.GetInvoiceByChecksum(checksum)
				if err != nil {
					return err
				}
				if err = pdf.PrintLocalInvoice(printingEngine, registryMeta, filename); err != nil {
					return err
				}
			}
		}

		finished = !resp.HasMore

		page += 1
	}

	return registry.Save("")
}
