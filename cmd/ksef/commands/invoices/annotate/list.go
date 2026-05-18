package annotate

import (
	annotationspkg "ksef/internal/invoicesdb/annotations"
	"ksef/internal/invoicesdb/shared"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list [faktura.xml]",
	Short:   "wyświetla adnotacje i flagi JPK z faktury",
	RunE:    runList,
	Args:    cobra.ExactArgs(1),
	PreRunE: initAnnotationRule,
}

func runList(cmd *cobra.Command, args []string) error {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "wiersz"},
			{Align: simpletable.AlignCenter, Text: "nazwa"},
			{Align: simpletable.AlignCenter, Text: "stawka VAT"},
			{Align: simpletable.AlignCenter, Text: "uwagi"},
		},
	}

	for _, item := range ctx.XMLInvoice.Items {
		itemHash := item.Hash()

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
				Text: getItemAnnotation(ctx.Invoice, itemHash, ctx.AnnotationsMgr),
			},
		})
	}

	table.SetStyle(simpletable.StyleDefault)
	table.Println()

	return nil
}

func getItemAnnotation(invoice *monthlyregistry.Invoice, hash shared.ItemHash, mgr *annotationspkg.Annotations) string {
	if rule := mgr.GetItemRule(invoice, hash); rule != nil {
		return rule.String()
	}
	return ""
}
