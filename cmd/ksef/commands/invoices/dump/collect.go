package dump

import (
	"fmt"
	"os"
	"ksef/internal/invoicesdb/annotations"
	"ksef/internal/invoicesdb/shared"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
)

// InvoiceAnnotations holds annotations for a single invoice's items.
type InvoiceAnnotations struct {
	Invoice    *monthlyregistry.Invoice
	RefNo      string
	SellerName string
	ItemRules  []ItemAnnotation
}

// ItemAnnotation holds the annotation for a single invoice item.
type ItemAnnotation struct {
	Item       monthlyregistry.Item
	Annotation *shared.Annotation
}

// InvoiceFilePaths holds the paths for a single invoice's PDF and XML files.
type InvoiceFilePaths struct {
	PDF string
	XML string
}

// CollectedInvoices holds all invoices and their annotations for a month.
type CollectedInvoices struct {
	InvoiceAnnotations []InvoiceAnnotations
	FilePaths          []InvoiceFilePaths
}

// CollectInvoices scans the registry for eligible invoices, collects their file
// paths and annotations. This is the shared logic used by both ZIP and PDF dump
// commands.
func CollectInvoices(
	registry *monthlyregistry.Registry,
	annotationsMgr *annotations.Annotations,
) (*CollectedInvoices, error) {
	collected := &CollectedInvoices{}

	for _, invoice := range registry.JPKEligibleInvoices() {
		fileNames := registry.InvoiceFilename(invoice)

		if _, err := os.Stat(fileNames.PDF); err != nil {
			return nil, fmt.Errorf("plik PDF nie istnieje: %s (faktura %s)", fileNames.PDF, invoice.RefNo)
		}

		if _, err := os.Stat(fileNames.XML); err != nil {
			return nil, fmt.Errorf("plik XML nie istnieje: %s (faktura %s)", fileNames.XML, invoice.RefNo)
		}

		collected.FilePaths = append(collected.FilePaths, InvoiceFilePaths{PDF: fileNames.PDF, XML: fileNames.XML})

		if invoiceAnnotations := collectInvoiceAnnotations(fileNames, invoice, annotationsMgr); invoiceAnnotations != nil {
			collected.InvoiceAnnotations = append(collected.InvoiceAnnotations, *invoiceAnnotations)
		}
	}

	return collected, nil
}

// collectInvoiceAnnotations parses the invoice's XML and collects annotations
// for each item. Annotations can be local (in the registry) or global (in
// subject settings). Only items that have at least one annotation are returned.
func collectInvoiceAnnotations(
	fileNames monthlyregistry.InvoiceFilename,
	invoice *monthlyregistry.Invoice,
	annotationsMgr *annotations.Annotations,
) *InvoiceAnnotations {
	xmlInvoice, _, err := monthlyregistry.ParseInvoice(fileNames.XML)
	if err != nil {
		return nil
	}

	var itemRules []ItemAnnotation
	for _, item := range xmlInvoice.Items {
		hash := item.Hash()
		if rule := annotationsMgr.GetItemRule(invoice, hash); rule != nil {
			itemRules = append(itemRules, ItemAnnotation{
				Item:       item,
				Annotation: rule,
			})
		}
	}

	if len(itemRules) == 0 {
		return nil
	}

	var sellerName string
	if invoice.Issuer != nil {
		sellerName = invoice.Issuer.Name
	}

	return &InvoiceAnnotations{
		Invoice:    invoice,
		RefNo:      invoice.RefNo,
		SellerName: sellerName,
		ItemRules:  itemRules,
	}
}
