package uploader

import (
	"bytes"
	"errors"
	"ksef/internal/invoice"
	"ksef/internal/logging"
	"ksef/internal/sei"
	"time"
)

const (
	invoiceTypeCorrection = "KOR"
)

var (
	errRetrievingContents = errors.New("unable to extract original invoice contents")
	errAddingItem         = errors.New("unable to add item")
)

func (u *Uploader) handleCorrection(i *sei.ParsedInvoice, originalInvoiceData *Invoice) error {
	currentInvoice := i.Invoice
	originalInvoice, err := originalInvoiceData.SourceDocument()
	if err != nil {
		return errors.Join(errRetrievingContents, err)
	}

	var correctionInvoice = &invoice.Invoice{}
	*correctionInvoice = *originalInvoice
	correctionInvoice.Type = invoiceTypeCorrection
	correctionInvoice.Correction = &invoice.CorrectionInfo{
		OriginalIssueDate: originalInvoice.Issued,
		KSeFRefNo:         originalInvoiceData.KSeFRefNo,
		RefNo:             originalInvoiceData.RefNo,
	}
	correctionInvoice.Clear()
	correctionInvoice.Attributes = originalInvoice.Attributes

	for itemNo, itemData := range currentInvoice.Items {
		oldItemData := originalInvoice.Items[itemNo]

		if oldItemData == itemData {
			logging.GenerateLogger.Debug("pozycje są identyczne", "item", itemNo)
			continue
		}

		var oldItemClone *invoice.InvoiceItem = &invoice.InvoiceItem{}
		*oldItemClone = *oldItemData

		oldItemClone.Before = true
		oldItemClone.RowNo = itemNo + 1

		itemData.RowNo = itemNo + 1

		if err = correctionInvoice.AddCorrectedItem(oldItemClone, itemData); err != nil {
			return errors.Join(errAddingItem, err)
		}
	}

	correctionInvoice.Issued = time.Now().Local()
	correctionInvoice.Number = "KOR/1/2025"
	i.Invoice = correctionInvoice

	var tmpBuffer bytes.Buffer
	i.ToXML(time.Time{}, &tmpBuffer)

	return nil
}
