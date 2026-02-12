package invoicesdb

import (
	"errors"
	"ksef/internal/invoice"
	"ksef/internal/logging"
	"reflect"
)

// disclaimer
// ----------
//
// this should probably be a single function with some fancy mathematics like picking a max(parsedInvoice.Items, originalInvoice.Items)
// but I was worried I'd make some easy mistake and / or the end code wouldn't be exacly clear.

func correction_RemoveAllItems(
	originalInvoice *invoice.Invoice,
	parsedInvoice *invoice.Invoice,
	correctionInvoice *invoice.Invoice,
) (err error) {
	for itemNo, originalItem := range originalInvoice.Items {
		var originalItemClone *invoice.InvoiceItem = &invoice.InvoiceItem{}
		*originalItemClone = *originalItem
		originalItemClone.Before = true
		originalItemClone.RowNo = itemNo + 1

		correctionItem := &invoice.InvoiceItem{RowNo: itemNo + 1}

		if err = correctionInvoice.AddCorrectedItem(originalItemClone, correctionItem); err != nil {
			return errors.Join(errAddingItem, err)
		}
	}

	return nil
}

func correction_ItemHasBeenAdded(
	originalInvoice *invoice.Invoice,
	parsedInvoice *invoice.Invoice,
	correctionInvoice *invoice.Invoice,
) (err error) {
	// item has been added therefore we first have to iterate through original items, compare them (if possible)
	// and then add the remaining item(s)
	// let's start with the comparison:

	for itemNo, originalItem := range originalInvoice.Items {
		correctedItem := parsedInvoice.Items[itemNo]

		if reflect.DeepEqual(correctedItem, originalItem) {
			logging.GenerateLogger.Debug("pozycje są identyczne", "item", itemNo+1)
			continue
		}

		var originalItemClone *invoice.InvoiceItem = &invoice.InvoiceItem{}
		*originalItemClone = *originalItem
		originalItemClone.Before = true
		originalItemClone.RowNo = itemNo + 1

		correctedItem.RowNo = itemNo + 1

		if err = correctionInvoice.AddCorrectedItem(originalItemClone, correctedItem); err != nil {
			return errors.Join(errAddingItem, err)
		}
	}

	// now we can iterate between len(originalInvoice.Items) and len(parsedInvoice.Items)
	// for instance - if the original invoice had 2 items and the corrected has 4
	// we now have to iterate between items 3 and 4 (indices 2 and 3, respectively)
	for itemNo := len(originalInvoice.Items); itemNo < len(parsedInvoice.Items); itemNo += 1 {
		// this effectively creates an empty row, thus marking it as non-existent
		// in the final XML
		originalItemClone := &invoice.InvoiceItem{
			Before: true,
			RowNo:  itemNo + 1,
		}

		correctedItem := parsedInvoice.Items[itemNo]
		correctedItem.RowNo = itemNo + 1

		if err = correctionInvoice.AddCorrectedItem(originalItemClone, correctedItem); err != nil {
			return errors.Join(errAddingItem, err)
		}
	}

	return nil
}

func correction_ItemHasBeenRemoved(
	originalInvoice *invoice.Invoice,
	parsedInvoice *invoice.Invoice,
	correctionInvoice *invoice.Invoice,
) (err error) {
	// this means that parsedInvoice has less items than the original invoice. therefore we first
	// have to iterate through both (stopping at last item from parsed invoice) to compare rows,
	// following by iterating through the leftovers and marking them as removed.
	// let's start with comparison (the easy part):
	for itemNo, correctedItem := range parsedInvoice.Items {
		originalItem := originalInvoice.Items[itemNo]

		if reflect.DeepEqual(correctedItem, originalItem) {
			logging.GenerateLogger.Debug("pozycje są identyczne", "item", itemNo+1)
			continue
		}

		var originalItemClone *invoice.InvoiceItem = &invoice.InvoiceItem{}
		*originalItemClone = *originalItem
		originalItemClone.Before = true
		originalItemClone.RowNo = itemNo + 1

		correctedItem.RowNo = itemNo + 1

		if err = correctionInvoice.AddCorrectedItem(originalItemClone, correctedItem); err != nil {
			return errors.Join(errAddingItem, err)
		}
	}

	// ok. now we need to iterate between len(parsedInvoice.Items) and len(originalInvoice.Items) and
	// mark all found rows as removed. For example - if the parsedInvoice has 2 items and the
	// originalInvoice had 4 - we have to iterate through items 3 and 4 (indices 2 and 3, respectively)
	for itemNo := len(parsedInvoice.Items); itemNo < len(originalInvoice.Items); itemNo += 1 {
		originalItemClone := &invoice.InvoiceItem{}
		*originalItemClone = *originalInvoice.Items[itemNo]
		originalItemClone.RowNo = itemNo + 1
		originalItemClone.Before = true

		// this effectively creates an empty row, thus marking it as non-existent
		// in the final XML
		correctedItem := &invoice.InvoiceItem{
			RowNo: itemNo + 1,
		}

		if err = correctionInvoice.AddCorrectedItem(originalItemClone, correctedItem); err != nil {
			return errors.Join(errAddingItem, err)
		}
	}

	return nil
}

func correction_NumberOfItemsEqual(
	originalInvoice *invoice.Invoice,
	parsedInvoice *invoice.Invoice,
	correctionInvoice *invoice.Invoice,
) (err error) {
	for itemNo, itemData := range parsedInvoice.Items {
		originalItem := originalInvoice.Items[itemNo]

		if reflect.DeepEqual(itemData, originalItem) {
			logging.GenerateLogger.Debug("pozycje są identyczne", "item", itemNo+1)
			continue
		}

		var originalItemClone *invoice.InvoiceItem = &invoice.InvoiceItem{}
		*originalItemClone = *originalItem

		originalItemClone.Before = true
		originalItemClone.RowNo = itemNo + 1

		itemData.RowNo = itemNo + 1

		if err = correctionInvoice.AddCorrectedItem(originalItemClone, itemData); err != nil {
			return errors.Join(errAddingItem, err)
		}
	}

	return nil
}
