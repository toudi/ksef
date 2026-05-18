package dump

import (
	"fmt"
	"ksef/cmd/ksef/flags"
	"ksef/internal/logging"
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mergePDFCommand = &cobra.Command{
	Use:     "pdf [year] [month]",
	Short:   "zrzuca faktury w formie PDF oraz notatki dla księgowych w jeden scalony plik PDF",
	Args:    cobra.MaximumNArgs(2),
	RunE:    generateDumpPDF,
	PreRunE: initializeDump,
}

func init() {
	pdfFlagSet := mergePDFCommand.Flags()
	flags.NIP(pdfFlagSet)
	pdfFlagSet.StringP(
		"output",
		"o",
		"",
		"Ścieżka pliku wyjściowego (domyślnie: invoices-merged-{year}-{month}.pdf)",
	)
}

// generateDumpPDF merges all PDF invoices and the accountant-notes PDF (if annotations exist)
// into a single PDF file.
func generateDumpPDF(cmd *cobra.Command, args []string) error {
	vip := viper.GetViper()

	// Determine output path
	outputPath := vip.GetString("output")
	if outputPath == "" {
		outputPath = fmt.Sprintf("invoices-merged-%d-%02d.pdf", dumpMonth.Year(), int(dumpMonth.Month()))
	}

	// Resolve output path to absolute
	absOutputPath, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("błąd podczas rozwiązywania ścieżki wyjściowej: %w", err)
	}

	// Collect all invoices, file paths, and annotations
	collected, err := CollectInvoices(dumpRegistry, dumpAnnotationsMgr)
	if err != nil {
		return err
	}

	if len(collected.FilePaths) == 0 {
		return fmt.Errorf("nie znaleziono żadnych plików PDF do połączenia")
	}

	// Build the list of PDFs to merge: annotations first (if any), then invoices
	allPDFs := []string{}
	if len(collected.InvoiceAnnotations) > 0 {
		workDir, err := os.MkdirTemp("", "invoices-merged-*")
		if err != nil {
			return fmt.Errorf("błąd tworzenia katalogu tymczasowego: %w", err)
		}
		defer os.RemoveAll(workDir)

		annotationsPDFPath := filepath.Join(workDir, "accountant-notes.pdf")
		if err := renderAccountantNotes(collected.InvoiceAnnotations, annotationsPDFPath); err != nil {
			return fmt.Errorf("nie udało się wygenerować accountant-notes.pdf: %w", err)
		}
		allPDFs = append(allPDFs, annotationsPDFPath)
	}

	for _, paths := range collected.FilePaths {
		allPDFs = append(allPDFs, paths.PDF)
	}

	// Merge the PDF files
	config := model.NewDefaultConfiguration()
	if err := api.MergeCreateFile(allPDFs, absOutputPath, false, config); err != nil {
		return fmt.Errorf("błąd podczas łączenia PDF-ów: %w", err)
	}

	logging.InvoicesDBLogger.Info("połączono faktury PDF", "plik", absOutputPath)
	return nil
}
