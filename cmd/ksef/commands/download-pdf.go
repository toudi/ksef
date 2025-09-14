package commands

import (
	"flag"
	"ksef/internal/pdf"
)

type downloadPDFCommand struct {
	Command
}
type downloadPDFArgsType struct {
	internalArgs pdf.DownloadPDFArgs
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
		},
	}

	DownloadPDFCommand.FlagSet.StringVar(
		&downloadPDFArgs.path,
		"p",
		"",
		"ścieżka do pliku rejestru",
	)
	DownloadPDFCommand.FlagSet.StringVar(
		&downloadPDFArgs.internalArgs.Output,
		"o",
		"",
		"ścieżka do zapisu PDF (domyślnie katalog pliku statusu + {nrRef}.pdf)",
	)
	DownloadPDFCommand.FlagSet.StringVar(
		&downloadPDFArgs.internalArgs.Invoice,
		"i",
		"",
		"numer faktury do pobrania. Wartość * oznacza pobranie wszystkich faktur z rejestru",
	)
	DownloadPDFCommand.FlagSet.StringVar(
		&downloadPDFArgs.internalArgs.IssuerToken,
		"token",
		"",
		"Token sesji interaktywnej lub nazwa zmiennej środowiskowej która go zawiera",
	)
	// DownloadPDFCommand.FlagSet.StringVar(&downloadPDFArgs.internalArgs.Token, "token", "", "token sesji")
	DownloadPDFCommand.FlagSet.BoolVar(
		&downloadPDFArgs.internalArgs.SaveXml,
		"xml",
		false,
		"zapisz źródłowy plik XML",
	)

	// API v2 zdaje się nie wspierać tej operacji więc do odwołania wyrejestrowuję komendę
	// registerCommand(&DownloadPDFCommand.Command)
}

func downloadPDFRun(c *Command) error {
	return ErrNotImplemented
	/*
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
			return pdf.DownloadPDFFromLocalFile(gateway, registry, &downloadPDFArgs.internalArgs)
		}

		return interactive.DownloadPDFFromAPI(gateway, &downloadPDFArgs.internalArgs, registry)
	*/
}
