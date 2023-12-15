package commands

import (
	"flag"
	"fmt"
	registryPkg "ksef/internal/registry"
	"ksef/internal/sei/api/client"
	"ksef/internal/sei/api/upo"
	"path/filepath"
)

type statusCommand struct {
	Command
}
type statusArgsType struct {
	path   string
	output string
	xml    bool
}

var StatusCommand *statusCommand
var statusArgs statusArgsType

func init() {
	StatusCommand = &statusCommand{
		Command: Command{
			Name:        "status",
			FlagSet:     flag.NewFlagSet("status", flag.ExitOnError),
			Description: "wysyła sprawdza status przesyłki i pobiera dokument UPO",
			Run:         statusRun,
			Args:        statusArgs,
		},
	}

	StatusCommand.FlagSet.StringVar(&statusArgs.path, "p", "", "ścieżka do pliku statusu")
	StatusCommand.FlagSet.StringVar(&statusArgs.output, "o", "", "ścieżka do zapisu UPO (domyślnie katalog pliku statusu + {nrRef}.pdf)")
	StatusCommand.FlagSet.BoolVar(&statusArgs.xml, "xml", false, "zapis UPO jako plik XML")

	registerCommand(&StatusCommand.Command)
}

func statusRun(c *Command) error {
	if statusArgs.path == "" {
		StatusCommand.FlagSet.Usage()
		return nil
	}

	registry, err := registryPkg.LoadRegistry(statusArgs.path)
	if err != nil {
		return fmt.Errorf("unable to load status from file: %v", err)
	}

	if registry.Environment == "" || registry.SessionID == "" {
		return fmt.Errorf("file deserialized correctly, but either environment or referenceNo are empty: %+v", registry)
	}

	gateway, err := client.APIClient_Init(registry.Environment)
	if err != nil {
		return fmt.Errorf("cannot initialize gateway: %v", err)
	}

	var outputFormat = upo.UPOFormatPDF
	if statusArgs.xml {
		outputFormat = upo.UPOFormatXML
	}

	if statusArgs.output == "" {
		statusArgs.output = filepath.Join(filepath.Dir(statusArgs.path), fmt.Sprintf("%s.%s", registry.SessionID, outputFormat))
	}

	if err = upo.DownloadUPO(gateway, registry, outputFormat, statusArgs.output); err != nil {
		return fmt.Errorf("unable to download UPO: %v", err)
	}

	return nil
}
