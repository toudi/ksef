package commands

import (
	"flag"
	"fmt"
	"ksef/internal/config"
	"ksef/internal/logging"
	registryPkg "ksef/internal/registry"
	v2 "ksef/internal/sei/api/client/v2"
	"ksef/internal/sei/api/session/interactive"
)

type uploadCommand struct {
	Command
}

type uploadArgsType struct {
	testGateway bool
	path        string
	interactive bool
	issuerToken string
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
			Args:        uploadArgs,
		},
	}

	UploadCommand.FlagSet.BoolVar(&uploadArgs.testGateway, "t", false, "użyj bramki testowej")
	UploadCommand.FlagSet.BoolVar(&uploadArgs.interactive, "i", false, "użyj sesji interaktywnej")
	UploadCommand.FlagSet.StringVar(
		&uploadArgs.issuerToken,
		"token",
		"",
		"Token sesji interaktywnej lub nazwa zmiennej środowiskowej która go zawiera",
	)
	UploadCommand.FlagSet.StringVar(
		&uploadArgs.path,
		"p",
		"",
		"ścieżka do katalogu z wygenerowanymi fakturami",
	)
	UploadCommand.FlagSet.BoolVar(
		&interactive.InteractiveSessionUploadParams.ForceUpload,
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

	var env config.APIEnvironment = config.APIEnvironmentProd
	if uploadArgs.testGateway {
		env = config.APIEnvironmentTest
	}

	registry, err := registryPkg.LoadRegistry(uploadArgs.path)
	if err != nil {
		return err
	}

	cli, err := v2.NewClient(c.Context, config.GetConfig(), env, v2.WithRegistry(registry))
	if err != nil {
		return fmt.Errorf("błąd inicjalizacji klienta: %v", err)
	}

	if uploadArgs.interactive {
		interactiveSession, err := cli.InteractiveSession()
		if err != nil {
			return err
		}
		if uploadArgs.issuerToken != "" {
			logging.AuthLogger.Warn("overriding KSeF token")
			if err = cli.Auth.SetKsefToken(uploadArgs.issuerToken); err != nil {
				logging.AuthLogger.Error("unable to override KSeF token")
				return err
			}
		}

		err = interactiveSession.UploadInvoices()
		if err == interactive.ErrProbablyUsedSend {
			fmt.Printf(
				"Wygląda na to, że poprzednio użyta została komenda 'upload' na tym rejestrze.\nJeśli na pewno chcesz ponowić wysyłkę, uzyj flagi '-f'\n",
			)
			return nil
		}
		return err
	}

	return fmt.Errorf("batch session not implemented yet")

	// batchSession := batch.BatchSessionInit(gateway)
	// return batchSession.UploadInvoices(uploadArgs.path)

}
