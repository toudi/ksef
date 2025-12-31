package monthlyregistry

import (
	"fmt"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/sei"

	"github.com/mozillazg/go-slugify"
)

var dirnameByType = map[InvoiceType]string{
	InvoiceTypeIssued:     "wystawione",
	InvoiceTypeReceived:   "otrzymane",
	InvoiceTypePayer:      "platnika",
	InvoiceTypeAuthorized: "strony-upowaznionej",
}

var ksefSubjectTypeToRegistryInvoiceType = map[invoices.SubjectType]InvoiceType{
	invoices.SubjectTypeRecipient:  InvoiceTypeReceived,
	invoices.SubjectTypePayer:      InvoiceTypePayer,
	invoices.SubjectTypeAuthorized: InvoiceTypeAuthorized,
}

func (r *Registry) getIssuedInvoiceFilename(invoiceNumber string, ordNo int) string {
	return fmt.Sprintf(
		"%s/%s/%04d-%s.xml",
		r.dir,
		dirnameByType[InvoiceTypeIssued],
		ordNo,
		slugify.Slugify(invoiceNumber),
	)
}

func (r *Registry) getInvoiceFilename(invoice *Invoice) string {
	if invoice.Type == InvoiceTypeIssued {
		return r.getIssuedInvoiceFilename(invoice.RefNo, invoice.OrdNum)
	}

	return fmt.Sprintf(
		"%s/%s/%04d-%s-%s.xml",
		r.dir,
		dirnameByType[invoice.Type],
		invoice.OrdNum,
		slugify.Slugify(invoice.Issuer.Name),
		slugify.Slugify(invoice.RefNo),
	)
}

func (r *Registry) GetDestFileName(inv *sei.ParsedInvoice, invoiceType InvoiceType) string {
	numInvoices := r.countInvoicesByType(invoiceType)

	if invoiceType == InvoiceTypeIssued {
		ordNo := numInvoices + 1
		// it's a slightly convoluted way of figuring out if we can reuse the filename
		// let's try to locate the invoice by the ref no. if it does not have
		// KSeFRefNo then we can reuse it's ord no and thus the filename
		existingInvoice, _ := r.getInvoiceByRefNo(inv.Invoice.Number)
		if existingInvoice != nil && existingInvoice.KSeFRefNo == "" {
			ordNo = existingInvoice.OrdNum
		}

		return r.getIssuedInvoiceFilename(inv.Invoice.Number, ordNo)
	}

	return fmt.Sprintf(
		"%s/%s/%04d-%s-%s.xml",
		r.dir,
		dirnameByType[invoiceType],
		numInvoices+1,
		slugify.Slugify(inv.Invoice.Issuer.Name),
		slugify.Slugify(inv.Invoice.Number),
	)
}

func (r *Registry) GetDestFileNameForAPIInvoice(subjectType invoices.SubjectType, inv invoices.InvoiceMetadata) string {
	invoiceType := ksefSubjectTypeToRegistryInvoiceType[subjectType]
	numInvoices := r.countInvoicesByType(invoiceType)

	return fmt.Sprintf(
		"%s/%s/%04d-%s-%s.xml",
		r.dir,
		dirnameByType[invoiceType],
		numInvoices+1,
		slugify.Slugify(inv.Seller.Name),
		slugify.Slugify(inv.InvoiceNumber),
	)
}

func (r *Registry) countInvoicesByType(invoiceType InvoiceType) (count int) {
	currentCount, exists := r.OrdNums[invoiceType]
	if !exists {
		return 0
	}
	return currentCount
}
