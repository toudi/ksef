package commands

import (
	"fmt"
	"ksef/cmd/ksef/commands/client"
	"ksef/cmd/ksef/flags"
	"ksef/internal/client/v2/session/interactive"
	"ksef/internal/config"
	registryPkg "ksef/internal/registry"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var uploadCommand = &cobra.Command{
	Use:   "upload",
	Short: "przesyła faktury do KSeF",
	RunE:  uploadRun,
}

type uploadArgsType struct {
	path                    string
	interactive             bool
	interactiveUploadParams interactive.UploadParams
}

var uploadArgs = &uploadArgsType{}

func init() {
	flags.AuthMethod(uploadCommand)
	var flags = uploadCommand.Flags()

	flags.BoolVarP(&uploadArgs.interactive, "interactive", "i", false, "użyj sesji interaktywnej")
	flags.StringVarP(&uploadArgs.path, "path", "p", "", "ścieżka do katalogu z wygenerowanymi fakturami")
	flags.BoolVarP(&uploadArgs.interactiveUploadParams.ForceUpload, "force", "f", false, "potwierdź wysyłkę faktur pomimo istniejących sum kontrolnych")

	_ = cobra.MarkFlagRequired(flags, "path")
	flags.SortFlags = false
}

func uploadRun(cmd *cobra.Command, _ []string) error {
	var ctx = cmd.Context()
	var vip = viper.GetViper()

	registry, err := registryPkg.OpenOrCreate(uploadArgs.path)
	if err != nil {
		return err
	}
	if registry.Environment == "" {
		registry.Environment = config.GetGateway(vip)
	}

	defer registry.Save("")

	// load invoice collection to retrieve the issuer
	collection, err := registry.InvoiceCollection()
	if err != nil {
		return err
	}

	if registry.Issuer == "" {
		registry.Issuer = collection.Issuer
	}

	cli, err := client.InitClient(cmd)
	if err != nil {
		return fmt.Errorf("błąd inicjalizacji klienta: %v", err)
	}

	defer cli.Logout()

	if uploadArgs.interactive {
		interactiveSession, err := cli.InteractiveSession()
		if err != nil {
			return err
		}

		err = interactiveSession.UploadInvoices(ctx, uploadArgs.interactiveUploadParams)
		if err == interactive.ErrProbablyUsedSend {
			fmt.Printf(
				"Wygląda na to, że poprzednio użyta została komenda 'upload' na tym rejestrze.\nJeśli na pewno chcesz ponowić wysyłkę, uzyj flagi '-f'\n",
			)
			return nil
		}
		return err
	}

	batchSession, err := cli.BatchSession()
	if err != nil {
		return err
	}

	return batchSession.UploadInvoices(ctx)
}
