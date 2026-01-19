package commands

import (
	"encoding/xml"
	"errors"
	"ksef/internal/config"
	invoicesdbconfig "ksef/internal/invoicesdb/config"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/logging"
	"ksef/internal/pdf"
	pdfconfig "ksef/internal/pdf/config"
	"ksef/internal/runtime"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var renderPDFCommand = &cobra.Command{
	Use:     "pdf",
	Short:   "drukuje PDF dla wskazanego dokumentu",
	RunE:    renderPDF,
	Args:    cobra.ExactArgs(1),
	PreRunE: detectRuntimeProperties,
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

func detectRuntimeProperties(cmd *cobra.Command, args []string) error {
	vip := viper.GetViper()
	// let's grab the source filename
	sourceFilename := args[0]
	if !strings.HasSuffix(sourceFilename, ".xml") {
		return errors.New("to nie jest plik XML")
	}

	filenameParts := strings.Split(filepath.Dir(sourceFilename), string(filepath.Separator))
	// the dir looks like this:
	// data/<gateway>/<nip>/<year>/<month>/<source>
	// what we want to extract is:
	// data dir (-6'th index)
	// gateway  (-5'th index)
	// nip      (-4'th index)
	nip := filenameParts[len(filenameParts)-4]
	environmentId := filenameParts[len(filenameParts)-5]
	dataDir := filenameParts[len(filenameParts)-6]
	runtime.SetNIP(vip, nip)
	runtime.SetEnvironment(vip, environmentId)
	invoicesdbconfig.SetDataDir(vip, dataDir)
	return nil
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
	engineConfig, err := pdfConfig.GetEngine(pdfconfig.UsageSelector{
		Usage: "upo",
	})
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
	engineConfig, err := pdfConfig.GetEngine(pdfconfig.UsageSelector{
		Usage:        invoiceMeta.Usage,
		Participants: invoiceMeta.Participants,
	})
	if err != nil {
		return err
	}
	logging.PDFRendererLogger.Debug("selected engine", "engine", engineConfig)
	printer, err := pdf.GetEngine(engineConfig)
	if err != nil {
		return err
	}
	return printer.PrintInvoice(invoiceXML, output, invoiceMeta)
}
