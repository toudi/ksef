package sei

import (
	"io"
	"ksef/internal/invoice"
	"time"
)

type ParsedInvoice struct {
	Invoice *invoice.Invoice
	sei     *SEI
}

func (p *ParsedInvoice) ToXML(generationTime time.Time, dst io.Writer) error {
	// ok so basically with this trick we can temporarily override the generation time.
	// this is important because when we're re-parsing the same file we need to
	// understand if we've already sent the invoice over to KSeF or not. But this is
	// tricky due to corrections. Henceforth we need to calculate the checksum of the
	// invoice *given the generation time*
	// let's go through this step by step:
	// 1. Initialize generation time to current time. This would be then propagated in
	// DataWytworzeniaFa field.
	p.Invoice.GenerationTime = time.Now().Local()
	// if generation time is given this means that we are trying to re-create the invoice
	// to check if we have already sent it to KSeF. As stated before, timestamp is part
	// of the checksum.
	if !generationTime.IsZero() {
		// override the one in parsed invoice with the one that we received from annual registry
		// this means - a timestamp which is the last known generation timestamp for this
		// invoice.
		p.Invoice.GenerationTime = generationTime
	}

	rootNode, err := p.sei.generator.InvoiceToXMLTree(p.Invoice)
	if err != nil {
		return err
	}

	if _, err = dst.Write([]byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")); err != nil {
		return err
	}

	if err = rootNode.DumpToWriter(dst, 0); err != nil {
		return err
	}

	return nil
}

func (s *SEI) invoiceReady(i *invoice.Invoice) error {
	return s.invoiceReadyFunc(&ParsedInvoice{Invoice: i, sei: s})
}
