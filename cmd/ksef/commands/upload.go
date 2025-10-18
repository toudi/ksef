package commands

import (
	"fmt"
	"ksef/cmd/ksef/flags"
	"ksef/internal/config"
	environmentPkg "ksef/internal/environment"
	registryPkg "ksef/internal/registry"
	v2 "ksef/internal/sei/api/client/v2"
	"ksef/internal/sei/api/client/v2/session/interactive"

	"github.com/spf13/cobra"
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
	env = environmentPkg.FromContext(ctx)

	registry, err := registryPkg.OpenOrCreate(uploadArgs.path)
	if err != nil {
		return err
	}
	if registry.Environment == "" {
		registry.Environment = env
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

	authValidator, err := authChallengeValidatorInstance(cmd, collection.Issuer, env)
	if err != nil {
		return err
	}

	cli, err := v2.NewClient(
		ctx,
		config.GetConfig(),
		env, v2.WithRegistry(registry), v2.WithAuthValidator(authValidator),
	)
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
