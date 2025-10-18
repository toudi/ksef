package commands

import (
	"fmt"
	"ksef/cmd/ksef/flags"
	"ksef/internal/config"
	registryPkg "ksef/internal/registry"
	v2 "ksef/internal/sei/api/client/v2"
	"ksef/internal/sei/api/client/v2/upo"

	"github.com/spf13/cobra"
)

const (
	flagNameRegistry = "registry"
)

var statusCommand = &cobra.Command{
	Use:   "status",
	Short: "sprawdza status wysłanych faktur i pobiera dokument UPO",
	RunE:  statusRun,
}

var registryPath string

var upoDownloaderParams upo.UPODownloaderParams

func init() {
	flagSet := statusCommand.Flags()
	flags.AuthMethod(statusCommand)
	cobra.MarkFlagRequired(flagSet, flagNameRegistry)

	upoDownloaderParams.Format = upo.UPOFormatPDF

	flagSet.StringVarP(&registryPath, flagNameRegistry, "r", "", "ścieżka do katalogu z rejestrem")
	flagSet.StringVarP(&upoDownloaderParams.Path, "output", "o", "", "ścieżka do zapisu UPO (domyślnie katalog rejestru + {nrRef}.pdf)")
	flagSet.BoolVarP(&upoDownloaderParams.Mkdir, "m", "", false, "stwórz ktalog do zapisu, jeśli wskazany nie istnieje")
	flagSet.BoolFunc("xml", "zapis UPO jako plik XML", func(s string) error {
		upoDownloaderParams.Format = upo.UPOFormatXML
		return nil
	})

	flagSet.SortFlags = false
}

func statusRun(cmd *cobra.Command, _ []string) error {
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

	env = registry.Environment

	authValidator, err := authChallengeValidatorInstance(cmd, registry.Issuer, env)
	if err != nil {
		return err
	}

	cli, err := v2.NewClient(cmd.Context(), config.GetConfig(), registry.Environment, v2.WithRegistry(registry), v2.WithAuthValidator(authValidator))
	if err != nil {
		return fmt.Errorf("błąd inicjalizacji klienta: %v", err)
	}

	defer cli.Logout()

	if upoDownloaderParams.Path == "" {
		upoDownloaderParams.Path = registry.Dir
	}

	return cli.UploadSessionsStatusCheck(cmd.Context(), upoDownloaderParams)
}
