package recorder

import (
	"bytes"
	"ksef/internal/invoice"
	"ksef/internal/sei"
	"time"
)

type parsedInvoiceRecorder struct {
	Invoices           []*invoice.Invoice
	XMLInvoices        []string
	lastGenerationTime time.Time // for deterministic content checking
	contentBuffer      bytes.Buffer
}

func (r *parsedInvoiceRecorder) Ready(i *sei.ParsedInvoice) error {
	r.Invoices = append(r.Invoices, i.Invoice)
	if err := i.ToXML(r.lastGenerationTime, &r.contentBuffer); err != nil {
		return err
	}
	r.XMLInvoices = append(r.XMLInvoices, r.contentBuffer.String())
	r.contentBuffer.Reset()
	r.lastGenerationTime = r.lastGenerationTime.Add(time.Second)
	return nil
}

func NewRecorder() *parsedInvoiceRecorder {
	return &parsedInvoiceRecorder{
		lastGenerationTime: time.Date(2026, 4, 14, 16, 0, 0, 0, time.UTC),
	}
}
