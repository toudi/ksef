package invoices

import (
	"context"
	"fmt"
	"ksef/internal/http"
	"ksef/internal/registry"
	"ksef/internal/sei/api/client/v2/types/invoices"
	"ksef/internal/utils"
	baseHTTP "net/http"
	"os"
	"path"
)

const (
	endpointDownloadInvoice = "/api/v2/invoices/ksef/%s"
)

type invoiceDownloader struct {
	httpClient *http.Client
	registry   *registry.InvoiceRegistry
	targetPath string
}

func NewInvoiceDownloader(httpClient *http.Client, registry *registry.InvoiceRegistry) *invoiceDownloader {
	return &invoiceDownloader{
		httpClient: httpClient,
		registry:   registry,
		targetPath: registry.Dir,
	}
}

func (d *invoiceDownloader) Download(
	ctx context.Context,
	invoiceMeta invoices.InvoiceMetadata,
	subjectType invoices.SubjectType,
) (outputFilename string, checksum string, err error) {
	outputFilename = d.registry.GetTargetFilename(invoiceMeta, subjectType)

	if err = os.MkdirAll(path.Dir(outputFilename), 0755); err != nil {
		return "", "", err
	}

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return "", "", err
	}
	defer outputFile.Close()

	_, err = d.httpClient.Request(
		ctx, http.RequestConfig{
			DestWriter:     outputFile,
			ExpectedStatus: baseHTTP.StatusOK,
		}, fmt.Sprintf(endpointDownloadInvoice, invoiceMeta.KSeFNumber),
	)

	if err != nil {
		return "", "", err
	}

	fileMeta, err := utils.FileSizeAndSha256Hash(outputFilename)
	if err != nil {
		return "", "", err
	}

	d.registry.AddInvoice(invoiceMeta, checksum)

	return outputFilename, fileMeta.Hash, nil
}
