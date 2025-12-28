package monthlyregistry

import "strings"

type InvoiceFilename struct {
	XML string
	PDF string
}

func (r *Registry) InvoiceFilename(i *Invoice) InvoiceFilename {
	var sourceFilename string = r.getInvoiceFilename(i)

	return InvoiceFilename{
		XML: sourceFilename,
		PDF: strings.Replace(sourceFilename, ".xml", ".pdf", 1),
	}
}
