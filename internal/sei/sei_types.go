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
	p.Invoice.GenerationTime = time.Now().Local()
	// ok so basically with this trick we can temporarily override the generation time.
	// this is important because when we're re-parsing the same file we need to
	// understand if we've already sent the invoice over to KSeF or not. But this is
	// tricky due to corrections. Henceforth we need to calculate the checksum of the
	// invoice *given the generation time*
	if !generationTime.IsZero() {
		origGenerationTime := p.Invoice.GenerationTime
		p.Invoice.GenerationTime = generationTime
		defer func() {
			p.Invoice.GenerationTime = origGenerationTime
		}()
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
