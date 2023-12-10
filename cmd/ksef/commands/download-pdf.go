package commands

import (
	"flag"
	"fmt"
	"ksef/internal/sei/api/client"
	"ksef/internal/sei/api/pdf"
	"ksef/internal/sei/api/status"
	"path/filepath"
)

type downloadPDFCommand struct {
	Command
}
type downloadPDFArgsType struct {
	path      string
	output    string
	invoiceNo string
}

var DownloadPDFCommand *downloadPDFCommand
var downloadPDFArgs downloadPDFArgsType

func init() {
	DownloadPDFCommand = &downloadPDFCommand{
		Command: Command{
			Name:        "download-pdf",
			FlagSet:     flag.NewFlagSet("download-pdf", flag.ExitOnError),
			Description: "pobiera wizualizację PDF dla wskazanej faktury",
			Run:         downloadPDFRun,
			Args:        downloadPDFArgs,
		},
	}

	DownloadPDFCommand.FlagSet.StringVar(&downloadPDFArgs.path, "p", "", "ścieżka do pliku statusu")
	DownloadPDFCommand.FlagSet.StringVar(&downloadPDFArgs.output, "o", "", "ścieżka do zapisu PDF (domyślnie katalog pliku statusu + {nrRef}.pdf)")
	DownloadPDFCommand.FlagSet.StringVar(&downloadPDFArgs.invoiceNo, "i", "", "numer faktury do pobrania")

	registerCommand(&DownloadPDFCommand.Command)
}

func downloadPDFRun(c *Command) error {
	if downloadPDFArgs.path == "" || downloadPDFArgs.invoiceNo == "" {
		DownloadPDFCommand.FlagSet.Usage()
		return nil
	}

	statusInfo, err := status.StatusFromFile(downloadPDFArgs.path)
	if err != nil {
		return fmt.Errorf("unable to load status from file: %v", err)
	}

	if statusInfo.Environment == "" || statusInfo.SessionID == "" {
		return fmt.Errorf("file deserialized correctly, but either environment or referenceNo are empty: %+v", statusInfo)
	}

	gateway, err := client.APIClient_Init(statusInfo.Environment)
	if err != nil {
		return fmt.Errorf("cannot initialize gateway: %v", err)
	}

	if downloadPDFArgs.output == "" {
		downloadPDFArgs.output = filepath.Dir(downloadPDFArgs.path)
	}

	if err = pdf.DownloadPDF(gateway, statusInfo, downloadPDFArgs.invoiceNo, downloadPDFArgs.output); err != nil {
		return fmt.Errorf("unable to download PDF: %v", err)
	}

	return nil
}
