package jpk

import (
	"errors"
	"fmt"
	"ksef/internal/invoicesdb/jpk"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/invoicesdb/shared"

	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var jpkDumpItems = &cobra.Command{
	Use:     "dump",
	Short:   "wyświetla tabelę pozycji z faktury kosztowej",
	Args:    cobra.ExactArgs(1),
	RunE:    dumpInvoiceItems,
	PreRunE: initJPKManagerFromInvoiceFile,
}

var (
	errNotACostInvoice = errors.New("to nie jest faktura kosztowa")
	errUnknownInvoice  = errors.New("nie odnaleziono faktury w rejestrze")
)

func initJPKManagerFromInvoiceFile(cmd *cobra.Command, args []string) error {
	var err error
	vip := viper.GetViper()

	invoiceRegistry, err = monthlyregistry.OpenFromInvoiceFile(args[0])
	if err != nil {
		return err
	}

	jpkManager, err = jpk.Manager(
		vip,
		jpk.WithMonthlyRegistry(invoiceRegistry),
	)
	if err != nil {
		return err
	}

	xmlInvoice, invoiceChecksum, err = monthlyregistry.ParseInvoice(args[0])
	if err != nil {
		return err
	}

	invoice = invoiceRegistry.GetInvoiceByChecksum(invoiceChecksum)

	if invoice == nil {
		return errUnknownInvoice
	}

	if invoice.Type != monthlyregistry.InvoiceTypeReceived {
		return errNotACostInvoice
	}

	return nil
}

func dumpInvoiceItems(cmd *cobra.Command, args []string) error {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "wiersz"},
			{Align: simpletable.AlignCenter, Text: "nazwa"},
			{Align: simpletable.AlignCenter, Text: "stawka VAT"},
			{Align: simpletable.AlignCenter, Text: "wyłącz z raportu"},
			{Align: simpletable.AlignCenter, Text: "raportuj VAT 50%"},
		},
	}

	for _, item := range xmlInvoice.Items {
		itemHash := shared.ItemHash{
			Name: item.Name,
		}

		table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
			{
				Text: item.OrdNo,
			},
			{
				Text: item.Name,
			},
			{
				Text: item.VATRate,
			},
			{
				Text: fmt.Sprintf("%t", jpkManager.ItemIsExcluded(invoice, itemHash)),
			},
			{
				Text: fmt.Sprintf("%t", jpkManager.ItemHasVat50PercentFlag(invoice, itemHash)),
			},
		})
	}

	table.SetStyle(simpletable.StyleDefault)
	table.Println()

	return nil
}
