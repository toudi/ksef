package invoices

import (
	"bytes"
	"encoding/json"
	"errors"
	"ksef/cmd/ksef/flags"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/runtime"
	"os"
	"text/template"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ksefMetaCommand = &cobra.Command{
	Short: "Zwraca metadane KSeF faktury",
	Use:   "ksef-meta [numer-faktury]",
	Args:  cobra.ExactArgs(1),
	RunE:  getKsefMeta,
}

var errInvoiceNotFound = errors.New("Nie znaleziono faktury")

func init() {
	flagSet := ksefMetaCommand.Flags()
	flagSet.VarP(flags.StringChoice([]string{
		"json",
		"yaml",
		"toml",
	}), "format", "f", "format wyjścia (domyślnie: tekstowy)")
	flagSet.StringP("output", "o", "", "plik wyjścia")
	flags.NIP(flagSet)

	InvoicesCommand.AddCommand(ksefMetaCommand)
}

type QRCodes struct {
	Invoice string `yaml:"invoice" toml:"invoice" json:"invoice"`
	Offline string `yaml:"offline,omitempty" toml:"offline,omitempty" json:"offline,omitempty"`
}
type ksefInvoiceMeta struct {
	RefNo     string  `yaml:"ref-no" toml:"ref-no" json:"ref-no"`
	KSeFRefNo string  `yaml:"ksef-ref-no" toml:"ksef-ref-no" json:"ksef-ref-no"`
	QRCodes   QRCodes `yaml:"qrcodes" toml:"qrcodes" json:"qrcodes"`
}

func getKsefMeta(cmd *cobra.Command, args []string) error {
	vip := viper.GetViper()
	if err := runtime.CheckNIPIsSet(vip); err != nil {
		return err
	}

	refNo := args[0]

	month := time.Now()
	month = time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.Local)

	var invoice *monthlyregistry.Invoice

	// let's start from the most recent month and iterate back
	// until we hit a match.
	for range 12 {
		reg, err := monthlyregistry.OpenForMonth(vip, month)
		if err != nil {
			month = month.AddDate(0, -1, 0)
			continue
		}
		if invoice = reg.GetInvoice(func(i monthlyregistry.Invoice) bool {
			return i.RefNo == refNo && i.Type == monthlyregistry.InvoiceTypeIssued
		}); invoice != nil {
			break
		}

		month = month.AddDate(0, -1, 0)
	}

	if invoice == nil {
		return errInvoiceNotFound
	}

	return dumpInvoiceData(invoice, vip)
}

func dumpInvoiceData(invoice *monthlyregistry.Invoice, vip *viper.Viper) error {
	format := vip.GetString("format")
	output := vip.GetString("output")

	var content []byte
	var err error

	invoiceMeta := ksefInvoiceMeta{
		RefNo:     invoice.RefNo,
		KSeFRefNo: invoice.KSeFRefNo,
		QRCodes: QRCodes{
			Invoice: invoice.QRCodes.Invoice,
			Offline: invoice.QRCodes.Offline,
		},
	}

	switch format {
	case "toml":
		content, err = toml.Marshal(invoiceMeta)
	case "yaml":
		content, err = yaml.Marshal(invoiceMeta)
	case "json":
		content, err = json.Marshal(invoiceMeta)
	default:
		content, err = dumpInvoiceMeta(invoiceMeta)
	}

	if err != nil {
		return err
	}

	outputWriter := os.Stdout

	if output != "" && output != "-" {
		outputWriter, err = os.Create(output)
		if err != nil {
			return err
		}
		defer outputWriter.Close()
	}

	_, err = outputWriter.Write(content)

	return err
}

const invoiceMetaTextTemplate = `
Numer faktury       : {{.RefNo}}
Numer faktury w KSeF: {{.KSeFRefNo}}
Kody QR
  Weryfikacyjny     : {{.QRCodes.Invoice}}
{{- if .QRCodes.Offline }}
  Offline           : {{.QRCodes.Offline}}
{{- end }}
`

func dumpInvoiceMeta(meta ksefInvoiceMeta) ([]byte, error) {
	var buffer bytes.Buffer

	tmpl, _ := template.New("").Parse(invoiceMetaTextTemplate)
	err := tmpl.Execute(&buffer, meta)

	return bytes.TrimLeft(buffer.Bytes(), "\n"), err
}
