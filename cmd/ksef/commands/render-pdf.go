package commands

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"ksef/internal/config"
	"ksef/internal/pdf"
	registryPkg "ksef/internal/registry"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var renderPDFCommand = &cobra.Command{
	Use:   "pdf",
	Short: "drukuje PDF dla wskazanego dokumentu",
}

var renderInvoicePDFCommand = &cobra.Command{
	Use:   "invoice",
	Short: "drukuje PDF dla wskazanej faktury",
	Args:  cobra.ExactArgs(1),
	RunE:  renderInvoicePDF,
}
var renderUPOPDFCommand = &cobra.Command{
	Use:   "upo",
	Short: "drukuje PDF dla wkazanego UPO",
	Args:  cobra.ExactArgs(1),
	RunE:  renderUPOPDF,
}

func init() {
	if err := config.PDFPrinterFlags(renderPDFCommand, renderPDFCommand.PersistentFlags()); err != nil {
		panic(err)
	}
	renderPDFCommand.AddCommand(renderInvoicePDFCommand)
	renderPDFCommand.AddCommand(renderUPOPDFCommand)
}

func renderInvoicePDF(cmd *cobra.Command, args []string) error {
	// let's grab the source filename
	invoiceSourceXML := args[0]

	// based on that, we can open registry
	registry, err := registryPkg.LoadRegistry(filepath.Dir(invoiceSourceXML))
	if err != nil {
		return err
	}

	// let's read the invoice into memory. we will need it anyway to pass it to the renderer
	// and to hash it (in order to generate the qrcode url)
	invoiceBytes, err := os.ReadFile(invoiceSourceXML)
	if err != nil {
		return err
	}
	// invoice, invoiceBytes, err := registryPkg.ReadAndParseInvoice(invoiceSourceXML)
	// if err != nil {
	// 	return err
	// }
	hash := sha256.Sum256(invoiceBytes)
	invoiceMeta, err := registry.GetInvoiceByChecksum(hex.EncodeToString(hash[:]))
	if err != nil {
		return err
	}

	// we're almost done. now let's prepare the qrcode URL
	// var qrcode string = "https://" + string(registry.Environment) + "/client-app/invoice/" + registry.Issuer + "/" + base64.URLEncoding.EncodeToString(hash[:])

	engine, err := pdf.GetLocalPrintingEngine()
	if err != nil {
		return err
	}
	basename, _ := strings.CutSuffix(filepath.Base(invoiceSourceXML), filepath.Ext(invoiceSourceXML))
	output := filepath.Join(registry.Dir, basename+".pdf")
	return engine.Print(base64.StdEncoding.EncodeToString(invoiceBytes), invoiceMeta, output)
}

func renderUPOPDF(cmd *cobra.Command, args []string) error {
	upoSourceXML := args[0]
	// based on that, we can open registry
	registry, err := registryPkg.LoadRegistry(filepath.Dir(upoSourceXML))
	if err != nil {
		return err
	}
	engine, err := pdf.GetLocalPrintingEngine()
	if err != nil {
		return err
	}
	upoBytes, err := os.ReadFile(upoSourceXML)
	if err != nil {
		return err
	}

	basename, _ := strings.CutSuffix(filepath.Base(upoSourceXML), filepath.Ext(upoSourceXML))
	output := filepath.Join(registry.Dir, basename+".pdf")

	return engine.PrintUPO(base64.StdEncoding.EncodeToString(upoBytes), output)
}
