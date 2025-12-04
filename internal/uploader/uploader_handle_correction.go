package uploader

import (
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
	ErrRecipientChanged   = errors.New("different recipient")
)

func (u *Uploader) handleCorrection(i *sei.ParsedInvoice, originalInvoiceData *Invoice) error {
	currentInvoice := i.Invoice
	originalInvoice, err := originalInvoiceData.SourceDocument()
	if err != nil {
		return errors.Join(errRetrievingContents, err)
	}

	var correctionInvoice = &invoice.Invoice{}
	*correctionInvoice = *originalInvoice
	correctionInvoice.Clear()
	correctionInvoice.Type = invoiceTypeCorrection
	correctionInvoice.Correction = &invoice.CorrectionInfo{
		OriginalIssueDate: originalInvoice.Issued,
		KSeFRefNo:         originalInvoiceData.KSeFRefNo,
		RefNo:             originalInvoiceData.RefNo,
	}
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
	correctionInvoice.Number = u.invoiceDB.getCorrectionNumber(correctionInvoice.Issued)
	i.Invoice = correctionInvoice

	if err = i.ToXML(time.Time{}, &u.contentBuffer); err != nil {
		return err
	}

	return nil
}
