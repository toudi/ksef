package commands

import (
	"flag"
	"fmt"
	registryPkg "ksef/internal/registry"
	"ksef/internal/sei/api/client"
	"ksef/internal/sei/api/upo"
	"ksef/internal/utils"
	"path/filepath"
)

type statusCommand struct {
	Command
}

type statusArgsType struct {
	path            string
	xml             bool
	downloadUPOArgs upo.DownloadUPOParams
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

	StatusCommand.FlagSet.StringVar(&statusArgs.path, "p", "", "ścieżka do pliku rejestru")
	StatusCommand.FlagSet.StringVar(
		&statusArgs.downloadUPOArgs.Output,
		"o",
		"",
		"ścieżka do zapisu UPO (domyślnie katalog pliku rejestru + {nrRef}.pdf)",
	)
	StatusCommand.FlagSet.BoolVar(
		&statusArgs.downloadUPOArgs.Mkdir,
		"m",
		false,
		"stwórz katalog, jeśli wskazany do zapisu nie istnieje",
	)
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
		return fmt.Errorf(
			"file deserialized correctly, but either environment or referenceNo are empty: %+v",
			registry,
		)
	}

	statusArgs.downloadUPOArgs.OutputFormat = upo.UPOFormatPDF

	if statusArgs.xml {
		statusArgs.downloadUPOArgs.OutputFormat = upo.UPOFormatXML
	}

	if statusArgs.downloadUPOArgs.Output == "" {
		statusArgs.downloadUPOArgs.Output = filepath.Dir(statusArgs.path)
	}

	outputPath, err := utils.ResolveFilepath(
		utils.FilepathResolverConfig{
			Path:  statusArgs.downloadUPOArgs.Output,
			Mkdir: statusArgs.downloadUPOArgs.Mkdir,
			DefaultFilename: fmt.Sprintf(
				"%s.%s",
				registry.SessionID,
				statusArgs.downloadUPOArgs.OutputFormat,
			),
		},
	)

	if err == utils.ErrDoesNotExistAndMkdirNotSpecified {
		return fmt.Errorf("wskazany katalog nie istnieje a nie użyłeś opcji `-m`")
	}
	if err != nil {
		return fmt.Errorf("błąd tworzenia katalogu wyjściowego: %v", err)
	}

	statusArgs.downloadUPOArgs.OutputPath = outputPath

	gateway, err := client.APIClient_Init(registry.Environment)
	if err != nil {
		return fmt.Errorf("cannot initialize gateway: %v", err)
	}

	if err = upo.DownloadUPO(gateway, registry, &statusArgs.downloadUPOArgs); err != nil {
		return fmt.Errorf("unable to download UPO: %v", err)
	}

	return nil
}
