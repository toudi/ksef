package commands

import (
	"flag"
	"fmt"
	"ksef/api"
)

type uploadCommand struct {
	Command
}

type uploadArgsType struct {
	testGateway bool
	path        string
	interactive bool
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
	UploadCommand.FlagSet.StringVar(&uploadArgs.path, "p", "", "ścieżka do katalogu z wygenerowanymi fakturami")

	registerCommand(&UploadCommand.Command)
}

func uploadRun(c *Command) error {
	if uploadArgs.path == "" {
		c.FlagSet.Usage()
		return nil
	}

	environment := api.ProductionEnvironment
	if uploadArgs.testGateway {
		environment = api.TestEnvironment
	}

	gateway, err := api.API_Init(environment)
	if err != nil {
		return fmt.Errorf("nieznane środowisko: %v", environment)
	}

	if uploadArgs.interactive {
		interactiveSession := gateway.InteractiveSessionInit()
		return interactiveSession.UploadInvoices(uploadArgs.path)
	}

	batchSession := gateway.BatchSessionInit()
	return batchSession.UploadInvoices(uploadArgs.path)

}
