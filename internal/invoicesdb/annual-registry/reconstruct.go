package annualregistry

import (
	"errors"
	"ksef/internal/invoice"
)

var (
	errNoOriginalContents = errors.New("no original invoice contents available for reconstruction")
	errReconstructItem    = errors.New("unable to reconstruct invoice item")
)

// ReconstructCurrentState applies all corrections, in chronological order, on
// top of the original invoice to produce the current state of the invoice in
// KSeF.
//
// Each correction is incremental: it contains pairs of items, Before=true (the
// item as it stood in the invoice state the correction was issued against) and
// Before=false (the corrected item). Although every correction must reference the
// original document, its content is a delta *relative* to the state produced by
// the previous correction, so each correction has to be applied in sequence -
// the last correction alone does not describe the current state.
func (i *Invoice) ReconstructCurrentState() (*invoice.Invoice, error) {
	if i.Contents == "" {
		return nil, errors.Join(errNoOriginalContents, errors.New("invoice RefNo: "+i.RefNo))
	}

	original, err := i.Unmarshall()
	if err != nil {
		return nil, err
	}

	// so we start with original items
	items := cloneItems(original.Items)

	for _, correction := range i.Corrections {
		corrInvoice, err := correction.Unmarshall()
		if err != nil {
			return nil, errors.Join(errReconstructItem, errors.New("correction: "+correction.RefNo), err)
		}
		// and apply each correction on top of the previous one to get the "current" state
		items = applyCorrection(items, corrInvoice.Items)
	}

	// Rebuild the invoice with the reconstructed items, using AddItem to
	// correctly calculate totals.
	result := &invoice.Invoice{}
	*result = *original
	result.Items = make([]*invoice.InvoiceItem, 0, len(items))
	result.TotalPerVATRate = make(map[string]invoice.Amount)
	result.Total = invoice.Amount{}

	for _, item := range items {
		if item == nil {
			continue
		}
		if err := result.AddItem(item); err != nil {
			return nil, errors.Join(errors.New("unable to add item during reconstruction"), err)
		}
	}

	return result, nil
}

// cloneItems clones the given items with their row numbers and before flags
// zeroed out, so that the reconstructed invoice items don't carry stale
// metadata from the source documents.
func cloneItems(items []*invoice.InvoiceItem) []*invoice.InvoiceItem {
	clones := make([]*invoice.InvoiceItem, len(items))
	for idx, item := range items {
		clone := *item
		clone.RowNo = 0
		clone.Before = false
		clones[idx] = &clone
	}
	return clones
}

// applyCorrection applies a single incremental correction on top of the
// current invoice state and returns the new state, repacked so that its rows
// are numbered 1..N (removed items dropped) - which is the numbering the next
// correction expects.
//
// The correction contains pairs of items: Before=true (skipped) and
// Before=false (applied as replacement at the row it references). An empty
// Before=false item signals that the referenced row was removed.
func applyCorrection(state []*invoice.InvoiceItem, corrItems []*invoice.InvoiceItem) []*invoice.InvoiceItem {
	// we do itemIdx+2 here to skip `before=true`. This works due to how we serialize items
	// in the registry - basically `before=true` always comes prior to rows with `before=false`
	for itemIdx := 0; itemIdx+1 < len(corrItems); itemIdx += 2 {
		beforeItem := corrItems[itemIdx]
		afterItem := corrItems[itemIdx+1]

		rowNo := afterItem.RowNo
		if rowNo == 0 {
			rowNo = beforeItem.RowNo
		}

		if afterItem.Description == "" && afterItem.Quantity.Amount == 0 && afterItem.UnitPrice.Amount == 0 {
			// Removal: nil out the item at this position
			if rowNo-1 < len(state) {
				state[rowNo-1] = nil
			}
		} else {
			clone := *afterItem
			clone.Before = false
			clone.RowNo = 0

			if rowNo-1 < len(state) {
				state[rowNo-1] = &clone
			} else {
				// The row goes beyond the current state (new item added):
				// pad with nils up to the new row, then append
				for len(state) < rowNo-1 {
					state = append(state, nil)
				}
				state = append(state, &clone)
			}
		}
	}

	// Repack: drop removed items so that the returned state is gap-free.
	repacked := make([]*invoice.InvoiceItem, 0, len(state))
	for _, item := range state {
		if item != nil {
			repacked = append(repacked, item)
		}
	}

	return repacked
}
