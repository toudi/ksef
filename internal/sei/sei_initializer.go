package sei

func WithInvoiceByRefNoExists(f func(refNo string) bool) func(s *SEI) {
	return func(s *SEI) {
		s.invoiceExistsFunc = f
	}
}

func WithInvoiceReadyFunc(f func(i *ParsedInvoice) error) func(s *SEI) {
	return func(s *SEI) {
		s.invoiceReadyFunc = f
	}
}
