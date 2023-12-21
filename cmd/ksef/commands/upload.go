package commands

import (
	"flag"
	"fmt"
	"ksef/internal/sei/api/client"
	"ksef/internal/sei/api/session/batch"
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
	UploadCommand.FlagSet.StringVar(&uploadArgs.issuerToken, "token", "", "Token sesji interaktywnej lub nazwa zmiennej środowiskowej która go zawiera")
	UploadCommand.FlagSet.StringVar(&uploadArgs.path, "p", "", "ścieżka do katalogu z wygenerowanymi fakturami")

	registerCommand(&UploadCommand.Command)
}

func uploadRun(c *Command) error {
	if uploadArgs.path == "" {
		c.FlagSet.Usage()
		return nil
	}

	environment := client.ProductionEnvironment
	if uploadArgs.testGateway {
		environment = client.TestEnvironment
	}

	gateway, err := client.APIClient_Init(environment)
	if err != nil {
		return fmt.Errorf("nieznane środowisko: %v", environment)
	}

	if uploadArgs.interactive {
		interactiveSession := interactive.InteractiveSessionInit(gateway)
		if uploadArgs.issuerToken != "" {
			interactiveSession.SetIssuerToken(uploadArgs.issuerToken)
		}
		return interactiveSession.UploadInvoices(uploadArgs.path)
	}

	batchSession := batch.BatchSessionInit(gateway)
	return batchSession.UploadInvoices(uploadArgs.path)

}
