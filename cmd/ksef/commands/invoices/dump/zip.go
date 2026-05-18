package dump

import (
	"archive/zip"
	"fmt"
	"ksef/cmd/ksef/flags"
	"ksef/internal/logging"
	"ksef/internal/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ZipDumpCommand = &cobra.Command{
	Use:     "dump [year] [month]",
	Short:   "eksportuje faktury z wybranego miesiąca do archiwum ZIP",
	Args:    cobra.MaximumNArgs(2),
	RunE:    generateDump,
	PreRunE: initializeDump,
}

func init() {
	flagSet := ZipDumpCommand.Flags()
	flags.NIP(flagSet)
	flagSet.StringP(
		"output",
		"o",
		"",
		"Ścieżka pliku wyjściowego (domyślnie: invoices-dump-{year}-{month}.zip)",
	)
	flagSet.BoolP(
		"xml",
		"x",
		false,
		"Uwzględnij pliki XML (domyślnie tylko PDF)",
	)

	ZipDumpCommand.AddCommand(mergePDFCommand)
}

func generateDump(cmd *cobra.Command, args []string) error {
	vip := viper.GetViper()

	// Determine output path
	outputPath := vip.GetString("output")
	if outputPath == "" {
		outputPath = fmt.Sprintf("invoices-dump-%d-%02d.zip", dumpMonth.Year(), int(dumpMonth.Month()))
	}

	includeXML := vip.GetBool("xml")

	// Collect all invoices, file paths, and annotations
	collected, err := CollectInvoices(dumpRegistry, dumpAnnotationsMgr)
	if err != nil {
		return err
	}

	// Create ZIP archive
	return createDumpArchive(collected, outputPath, includeXML)
}

func createDumpArchive(collected *CollectedInvoices, outputPath string, includeXML bool) error {
	zipFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add collected PDF and optional XML files to archive
	for _, paths := range collected.FilePaths {
		if err := addFileToArchive(zipWriter, paths.PDF); err != nil {
			return err
		}

		if includeXML {
			if err := addFileToArchive(zipWriter, paths.XML); err != nil {
				return err
			}
		}
	}

	// Render accountant-notes PDF if there are annotations and add it to the archive
	if len(collected.InvoiceAnnotations) > 0 {
		tmpPDF, err := os.CreateTemp("", "accountant-notes-*.pdf")
		if err != nil {
			return fmt.Errorf("błąd tworzenia pliku tymczasowego: %w", err)
		}
		tmpPath := tmpPDF.Name()
		tmpPDF.Close()
		defer os.Remove(tmpPath)

		if err := renderAccountantNotes(collected.InvoiceAnnotations, tmpPath); err != nil {
			logging.InvoicesDBLogger.Warn("nie udało się wygenerować accountant-notes.pdf", "err", err)
		} else if _, err := os.Stat(tmpPath); err == nil {
			zipEntry, err := zipWriter.Create("accountant-notes.pdf")
			if err != nil {
				return fmt.Errorf("błąd tworzenia wpisu PDF w archiwum: %w", err)
			}
			if err := utils.CopyFileToWriter(tmpPath, zipEntry); err != nil {
				return fmt.Errorf("błąd kopiowania accountant-notes.pdf do archiwum: %w", err)
			}
		}
	}

	logging.InvoicesDBLogger.Info("zarchiwizowano faktury", "plik", outputPath)
	return nil
}

func addFileToArchive(zipWriter *zip.Writer, filePath string) error {
	// Compute relative path for zip entry: split by separator and take last 2 parts
	// e.g. "/path/to/registry/wystawione/0001-invoice.xml" -> "wystawione/0001-invoice.xml"
	pathParts := strings.Split(filePath, string(filepath.Separator))
	if len(pathParts) < 2 {
		return fmt.Errorf("nieoczekiwana struktura ścieżki: %s", filePath)
	}
	relativePath := filepath.Join(pathParts[len(pathParts)-2:]...)

	// Create zip entry and copy file contents
	zipEntry, err := zipWriter.Create(relativePath)
	if err != nil {
		return fmt.Errorf("błąd podczas tworzenia wpisu w archiwum %s: %w", relativePath, err)
	}

	return utils.CopyFileToWriter(filePath, zipEntry)
}
