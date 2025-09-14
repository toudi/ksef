package commands

import (
	"flag"
	"fmt"
	"ksef/internal/config"
	registryPkg "ksef/internal/registry"
	v2 "ksef/internal/sei/api/client/v2"
	"ksef/internal/sei/api/client/v2/upo"
	"path/filepath"
)

type statusCommand struct {
	Command
}

var registryPath string

var StatusCommand *statusCommand
var upoDownloaderParams upo.UPODownloaderParams
var issuerToken string
var xml bool

func init() {
	upoDownloaderParams.Format = upo.UPOFormatPDF
	StatusCommand = &statusCommand{
		Command: Command{
			Name:        "status",
			FlagSet:     flag.NewFlagSet("status", flag.ExitOnError),
			Description: "wysyła sprawdza status przesyłki i pobiera dokument UPO",
			Run:         statusRun,
		},
	}

	flagSet := StatusCommand.FlagSet
	initAuthParams(flagSet)

	StatusCommand.FlagSet.StringVar(&registryPath, "p", "", "ścieżka do pliku rejestru")
	StatusCommand.FlagSet.StringVar(
		&upoDownloaderParams.Path,
		"o",
		"",
		"ścieżka do zapisu UPO (domyślnie katalog pliku rejestru + {nrRef}.pdf)",
	)
	StatusCommand.FlagSet.BoolVar(
		&upoDownloaderParams.Mkdir,
		"m",
		false,
		"stwórz katalog, jeśli wskazany do zapisu nie istnieje",
	)
	StatusCommand.FlagSet.BoolFunc("xml", "zapis UPO jako plik XML", func(s string) error {
		upoDownloaderParams.Format = upo.UPOFormatPDF
		return nil
	})

	registerCommand(&StatusCommand.Command)
}

func statusRun(c *Command) error {
	if registryPath == "" {
		StatusCommand.FlagSet.Usage()
		return nil
	}

	registry, err := registryPkg.LoadRegistry(registryPath)
	if err != nil {
		return fmt.Errorf("unable to load status from file: %v", err)
	}

	if registry.Environment == "" || registry.UploadSessions == nil {
		return fmt.Errorf(
			"file deserialized correctly, but either environment or referenceNo are empty: %+v",
			registry,
		)
	}

	authValidator := authValidatorInstance(registry.Issuer)

	cli, err := v2.NewClient(c.Context, config.GetConfig(), registry.Environment, v2.WithRegistry(registry), v2.WithAuthValidator(authValidator))
	if err != nil {
		return fmt.Errorf("błąd inicjalizacji klienta: %v", err)
	}

	if upoDownloaderParams.Path == "" {
		upoDownloaderParams.Path = filepath.Dir(registryPath)
	}

	return cli.UploadSessionsStatusCheck(c.Context, upoDownloaderParams)
}
