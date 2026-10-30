package metadata

import (
	"context"
	"ksef/core/downloader"
	"ksef/core/ksef"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/http"
	baseHttp "net/http"
	"strconv"

	downloadertypes "ksef/internal/client/v2/invoices/downloader/types"
	ratelimits "ksef/internal/client/v2/rate-limits"
)

const (
	endpointInvoicesMetadata = "/v2/invoices/query/metadata"
)

func InvoicesMetadataPage(
	ctx context.Context,
	httpClient *http.Client,
	subjectType ksef.SubjectType,
	params downloader.Config,
	pageNo int,
) (
	resp *invoices.InvoiceMetadataResponse, err error,
) {
	req := downloadertypes.InvoiceListRequest{
		SubjectType: subjectType,
		DateRange: downloadertypes.DateRange{
			DateType: downloadertypes.DateRangeStorage,
			From:     params.StartDate,
			To:       params.EndDate,
		},
	}

	_, err = httpClient.Request(
		ctx,
		http.RequestConfig{
			Method: baseHttp.MethodPost,
			QueryParams: map[string]string{
				"pageOffset": strconv.Itoa(pageNo),
				"pageSize":   strconv.Itoa(params.PageSize),
			},
			Body:            req,
			ContentType:     http.JSON,
			Dest:            &resp,
			DestContentType: http.JSON,
			ExpectedStatus:  baseHttp.StatusOK,
			OperationId:     ratelimits.OperationInvoiceMetadata,
		},
		endpointInvoicesMetadata,
	)

	return resp, err
}
