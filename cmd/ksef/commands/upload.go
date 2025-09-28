package commands

import (
	"flag"
	"fmt"
	"ksef/internal/config"
	registryPkg "ksef/internal/registry"
	v2 "ksef/internal/sei/api/client/v2"
	"ksef/internal/sei/api/client/v2/session/interactive"
)

type uploadCommand struct {
	Command
}

type uploadArgsType struct {
	path                    string
	interactive             bool
	interactiveUploadParams interactive.UploadParams
}

var UploadCommand *uploadCommand
var uploadArgs = &uploadArgsType{}

func init() {
	UploadCommand = &uploadCommand{
		Command: Command{
			Name:        "upload",
			FlagSet:     flag.NewFlagSet("upload", flag.ExitOnError),
			Description: "wysyła podpisany plik KSEF do bramki ministerstwa finansów",
			Run:         uploadRun,
		},
	}

	flagSet := UploadCommand.FlagSet
	initAuthParams(flagSet)
	testGatewayFlag(flagSet)

	UploadCommand.FlagSet.BoolVar(&uploadArgs.interactive, "i", false, "użyj sesji interaktywnej")
	UploadCommand.FlagSet.StringVar(
		&uploadArgs.path,
		"p",
		"",
		"ścieżka do katalogu z wygenerowanymi fakturami",
	)
	UploadCommand.FlagSet.BoolVar(
		&uploadArgs.interactiveUploadParams.ForceUpload,
		"f",
		false,
		"potwierdź wysyłkę faktur pomimo istniejących sum kontrolnych",
	)

	registerCommand(&UploadCommand.Command)
}

func uploadRun(c *Command) error {
	if uploadArgs.path == "" {
		c.FlagSet.Usage()
		return nil
	}

	registry, err := registryPkg.LoadRegistry(uploadArgs.path)
	if err != nil {
		return err
	}

	// load invoice collection to retrieve the issuer
	collection, err := registry.InvoiceCollection()
	if err != nil {
		return err
	}

	authValidator := authValidatorInstance(collection.Issuer)

	cli, err := v2.NewClient(c.Context, config.GetConfig(), environment, v2.WithRegistry(registry), v2.WithAuthValidator(authValidator))
	if err != nil {
		return fmt.Errorf("błąd inicjalizacji klienta: %v", err)
	}

	if uploadArgs.interactive {
		interactiveSession, err := cli.InteractiveSession()
		if err != nil {
			return err
		}

		err = interactiveSession.UploadInvoices(c.Context, uploadArgs.interactiveUploadParams)
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

	return batchSession.UploadInvoices(c.Context)
}
