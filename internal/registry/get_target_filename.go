package registry

import (
	"fmt"
	"ksef/internal/sei/api/client/v2/types/invoices"
	"path"
	"strings"

	"github.com/mozillazg/go-slugify"
)

var subjectTypeToDirname = map[invoices.SubjectType]string{
	invoices.SubjectTypeIssuer:     "wystawione",
	invoices.SubjectTypeRecipient:  "kosztowe",
	invoices.SubjectTypePayer:      "platnik",
	invoices.SubjectTypeAuthorized: "strona-upowazniona",
}

func (r *InvoiceRegistry) GetTargetFilename(
	invoice invoices.InvoiceMetadata,
	subjectType invoices.SubjectType,
) string {
	dir := r.GetDir()

	var sanitizedInvoiceNumber = slugify.Slugify(invoice.InvoiceNumber)
	var sanitizedFilename = strings.ToLower(
		fmt.Sprintf(
			"%03d-%s-%s.xml",
			len(r.Invoices)+1,
			invoice.Seller.Identifier,
			sanitizedInvoiceNumber,
		),
	)
	return path.Join(dir, subjectTypeToDirname[subjectType], sanitizedFilename)
}
