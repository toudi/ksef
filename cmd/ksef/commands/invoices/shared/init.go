package shared

import (
	"errors"
	"ksef/internal/invoicesdb/annotations"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	errNotACostInvoice = errors.New("to nie jest faktura zakupowa")
	errUnknownInvoice  = errors.New("nie odnaleziono faktury w rejestrze")
)

type InvoiceContext struct {
	AnnotationsMgr *annotations.Annotations
	Invoice        *monthlyregistry.Invoice
	XMLInvoice     *monthlyregistry.XMLInvoice
	InvoiceRegistry *monthlyregistry.Registry
}

func InitAnnotationsManagerFromInvoiceFile(cmd *cobra.Command, args []string) (*InvoiceContext, error) {
	var err error
	vip := viper.GetViper()

	invoiceRegistry, err := monthlyregistry.OpenFromInvoiceFile(args[0])
	if err != nil {
		return nil, err
	}

	annotationsMgr, err := annotations.Manager(
		vip,
		annotations.WithMonthlyRegistry(invoiceRegistry),
	)
	if err != nil {
		return nil, err
	}

	xmlInvoice, invoiceChecksum, err := monthlyregistry.ParseInvoice(args[0])
	if err != nil {
		return nil, err
	}

	invoice := invoiceRegistry.GetInvoiceByChecksum(invoiceChecksum)

	if invoice == nil {
		return nil, errUnknownInvoice
	}

	if invoice.Type != monthlyregistry.InvoiceTypeReceived {
		return nil, errNotACostInvoice
	}

	return &InvoiceContext{
		AnnotationsMgr: annotationsMgr,
		Invoice:        invoice,
		XMLInvoice:     xmlInvoice,
		InvoiceRegistry: invoiceRegistry,
	}, nil
}
