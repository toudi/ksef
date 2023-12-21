package commands

import (
	"flag"
	"fmt"
	"ksef/internal/registry"
	"ksef/internal/sei/api/client"
	"ksef/internal/sei/api/session/interactive"
	"path/filepath"
	"strings"
)

type downloadPDFCommand struct {
	Command
}
type downloadPDFArgsType struct {
	internalArgs registry.DownloadPDFArgs
	path         string
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

	DownloadPDFCommand.FlagSet.StringVar(&downloadPDFArgs.path, "p", "", "ścieżka do pliku rejestru")
	DownloadPDFCommand.FlagSet.StringVar(&downloadPDFArgs.internalArgs.Output, "o", "", "ścieżka do zapisu PDF (domyślnie katalog pliku statusu + {nrRef}.pdf)")
	DownloadPDFCommand.FlagSet.StringVar(&downloadPDFArgs.internalArgs.Invoice, "i", "", "numer faktury do pobrania. Wartość * oznacza pobranie wszystkich faktur z rejestru")
	DownloadPDFCommand.FlagSet.StringVar(&downloadPDFArgs.internalArgs.IssuerToken, "token", "", "Token sesji interaktywnej lub nazwa zmiennej środowiskowej która go zawiera")
	// DownloadPDFCommand.FlagSet.StringVar(&downloadPDFArgs.internalArgs.Token, "token", "", "token sesji")
	DownloadPDFCommand.FlagSet.BoolVar(&downloadPDFArgs.internalArgs.SaveXml, "xml", false, "zapisz źródłowy plik XML")

	registerCommand(&DownloadPDFCommand.Command)
}

func downloadPDFRun(c *Command) error {
	if downloadPDFArgs.path == "" || downloadPDFArgs.internalArgs.Invoice == "" {
		DownloadPDFCommand.FlagSet.Usage()
		return nil
	}

	registry, err := registry.LoadRegistry(downloadPDFArgs.path)
	if err != nil {
		return fmt.Errorf("unable to load registry from file: %v", err)
	}

	if registry.Environment == "" {
		return fmt.Errorf("file deserialized correctly, but environment is empty")
	}

	gateway, err := client.APIClient_Init(registry.Environment)
	if err != nil {
		return fmt.Errorf("cannot initialize gateway: %v", err)
	}

	if downloadPDFArgs.internalArgs.Output == "" {
		downloadPDFArgs.internalArgs.Output = filepath.Dir(downloadPDFArgs.path)
	}

	if strings.HasSuffix(strings.ToLower(downloadPDFArgs.internalArgs.Invoice), ".xml") {
		return registry.DownloadPDF(gateway, &downloadPDFArgs.internalArgs)
	}

	return interactive.DownloadPDFFromAPI(gateway, &downloadPDFArgs.internalArgs, registry)
}
