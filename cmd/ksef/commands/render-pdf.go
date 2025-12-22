package commands

import (
	"encoding/xml"
	"ksef/internal/config"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/pdf"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var renderPDFCommand = &cobra.Command{
	Use:   "pdf",
	Short: "drukuje PDF dla wskazanego dokumentu",
	RunE:  renderPDF,
}

const (
	flagNameOutput = "output"
)

type XMLFile struct {
	XMLName xml.Name
}

func init() {
	renderPDFCommand.Flags().StringP(flagNameOutput, "o", "", "plik wyjścia (jeśli go nie wskażesz, PDF zostanie utworzony w katalogu ze źródłowym XML)")
}

func renderPDF(cmd *cobra.Command, args []string) error {
	vip := viper.GetViper()

	// let's grab the source filename
	xmlContent, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}
	var xmlFile XMLFile
	if err = xml.Unmarshal(xmlContent, &xmlFile); err != nil {
		return err
	}

	pdfConfig, err := config.GetPDFPrinterConfig(vip)
	if err != nil {
		return err
	}

	output := strings.Replace(args[0], ".xml", ".pdf", 1)
	customOutput, err := cmd.Flags().GetString(flagNameOutput)
	if err != nil {
		return err
	}
	if customOutput != "" {
		output = customOutput
	}

	if xmlFile.XMLName.Local == "Potwierdzenie" {
		return renderUPO(pdfConfig, args[0], output)
	}
	return renderInvoice(pdfConfig, args[0], output)
}

func renderUPO(pdfConfig config.PDFPrinterConfig, upoXML string, output string) error {
	engineConfig, err := pdfConfig.GetEngine("upo")
	if err != nil {
		return err
	}
	printer, err := pdf.GetEngine(engineConfig)
	if err != nil {
		return err
	}

	return printer.PrintUPO(upoXML, output)
}

func renderInvoice(
	pdfConfig config.PDFPrinterConfig,
	invoiceXML string,
	output string,
) error {
	invoiceMeta, err := monthlyregistry.GetInvoicePrintingMeta(invoiceXML)
	if err != nil {
		return err
	}
	engineConfig, err := pdfConfig.GetEngine(invoiceMeta.Usage)
	if err != nil {
		return err
	}
	printer, err := pdf.GetEngine(engineConfig)
	if err != nil {
		return err
	}
	return printer.PrintInvoice(invoiceXML, output, invoiceMeta)
}
