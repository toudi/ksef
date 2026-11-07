package resources

import (
	"context"
	"fmt"
	"ksef/core/client"
	"ksef/core/invoicedb"
	"ksef/core/invoicedb/invoice"
	"ksef/internal/pdf"
	"ksef/internal/utils"
	"os"
	"path/filepath"

	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	pdfconfig "ksef/internal/pdf/config"

	"github.com/spf13/viper"
)

var _ client.InvoicePDFRenderer = (*pdfRenderer)(nil)

var invoiceTypeToPrinterUsage = map[invoice.InvoiceType]string{
	invoice.InvoiceTypeIssued:     "issued",
	invoice.InvoiceTypeReceived:   "received",
	invoice.InvoiceTypePayer:      "payer",
	invoice.InvoiceTypeAuthorized: "authorized",
}

type pdfRenderer struct {
	vip         *viper.Viper
	dataStorage invoicedb.InvoiceDataStorage
}

func NewPDFRenderer(
	vip *viper.Viper,
	dataStorage invoicedb.InvoiceDataStorage,
) client.InvoicePDFRenderer {
	return &pdfRenderer{
		vip:         vip,
		dataStorage: dataStorage,
	}
}

func (r *pdfRenderer) RenderInvoice(
	ctx context.Context,
	stored *invoicedb.StoredInvoice,
) error {
	usage := "invoice:" + invoiceTypeToPrinterUsage[stored.Invoice.Type]
	meta := &monthlyregistry.InvoicePrintingMeta{
		Usage: usage,
		Invoice: monthlyregistry.InvoiceMeta{
			KSeFRefNo: stored.Invoice.KSeFRefNo,
			QRCodes: monthlyregistry.InvoiceQRCodes{
				Invoice: stored.Invoice.QRCodes.Invoice,
				Offline: stored.Invoice.QRCodes.Offline,
			},
		},
		Printout: stored.Invoice.PrintoutData,
	}

	printer, err := pdf.GetInvoicePrinter(r.vip, pdfconfig.UsageSelector{Usage: usage})
	if err != nil {
		return err
	}

	// before we dive in more: I know it's hacky but here's the rationale behind this implementation.
	// if you're using typst renderer (which has access to local files) and you're also storing the
	// files locally (on disk) it makes no sense to abstract it and make the code always call
	// Store() (since typst can already write the file). but then we need to somehow instruct
	// the renderer if it can write directly. that's where this LocalPath comes in: it returns
	// a backend-computed path and a boolean that indicates if it is safe to write locally. This implies
	// that for backends that would be cloud-based (or even - database-based) the implementation could
	// simply return "", false and then the printer would use the .Store method (which, again, for
	// cloud-based storage implementations would probably mean issuing a HTTP PUT or something similar)
	xmlPath, xmlDirect := r.dataStorage.LocalPath(stored.StorageId(invoice.InvoiceFileXML))
	pdfPath, pdfDirect := r.dataStorage.LocalPath(stored.StorageId(invoice.InvoiceFilePDF))

	var workdir string
	if !xmlDirect || !pdfDirect {
		workdir, err = os.MkdirTemp("", "ksef-pdf-*")
		if err != nil {
			return err
		}
		defer os.RemoveAll(workdir)
	}

	if !xmlDirect {
		xmlPath = filepath.Join(workdir, "invoice.xml")
		src, err := r.dataStorage.Open(ctx, stored.StorageId(invoice.InvoiceFileXML))
		if err != nil {
			return err
		}
		defer src.Close()
		if err = utils.SaveReaderToFile(src, xmlPath); err != nil {
			return err
		}
	}

	if !pdfDirect {
		pdfPath = filepath.Join(workdir, "invoice.pdf")
	}

	if err := printer.PrintInvoice(xmlPath, pdfPath, meta); err != nil {
		return err
	}

	// as mentioned before - if we're using a filesystem-based storage (local) then typst can already write
	// to it directly and there's no need to use the abstract .Store() method
	if pdfDirect {
		return nil
	}

	pdfFile, err := os.Open(pdfPath)
	if err != nil {
		return err
	}
	defer pdfFile.Close()

	if _, err := r.dataStorage.Store(
		ctx,
		pdfFile,
		stored.StorageId(invoice.InvoiceFilePDF),
		false,
	); err != nil {
		return fmt.Errorf("unable to store rendered PDF: %w", err)
	}

	return nil
}
