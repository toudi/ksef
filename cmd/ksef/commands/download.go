package commands

import (
	"flag"
	"fmt"
	registryPkg "ksef/internal/registry"
	"ksef/internal/sei/api/client"
	"ksef/internal/sei/api/sync"
	"path"
	"time"
)

type syncInvoicesCommand struct {
	Command
}

type syncInvoicesArgsType struct {
	params      sync.SyncInvoicesConfig
	startDate   string
	testGateway bool
	refresh     string
}

var SyncInvoicesCommand *syncInvoicesCommand
var syncInvoicesArgs = &syncInvoicesArgsType{}

func init() {
	SyncInvoicesCommand = &syncInvoicesCommand{
		Command: Command{
			Name:        "download",
			FlagSet:     flag.NewFlagSet("download", flag.ExitOnError),
			Description: "Synchronizuje listę faktur z KSeF do katalogu lokalnego",
			Run:         syncInvoicesRun,
			Args:        syncInvoicesArgs,
		},
	}

	SyncInvoicesCommand.FlagSet.BoolVar(&syncInvoicesArgs.testGateway, "t", false, "użyj bramki testowej")
	SyncInvoicesCommand.FlagSet.StringVar(&syncInvoicesArgs.params.DestPath, "d", "", "Katalog docelowy")
	SyncInvoicesCommand.FlagSet.BoolVar(&syncInvoicesArgs.params.Income, "income", false, "Synchronizuj faktury przychodowe (Podmiot1)")
	SyncInvoicesCommand.FlagSet.BoolVar(&syncInvoicesArgs.params.Cost, "cost", false, "Synchronizuj faktury kosztowe (Podmiot2)")
	SyncInvoicesCommand.FlagSet.BoolVar(&syncInvoicesArgs.params.Subject3, "subject3", false, "Synchronizuj faktury podmiotu innego (Podmiot3)")
	SyncInvoicesCommand.FlagSet.BoolVar(&syncInvoicesArgs.params.SubjectAuthorized, "subjectAuthorized", false, "Synchronizuj faktury podmiotu upoważnionego (???)")
	SyncInvoicesCommand.FlagSet.StringVar(&syncInvoicesArgs.params.SubjectTIN, "nip", "", "Numer NIP podmiotu")
	SyncInvoicesCommand.FlagSet.StringVar(&syncInvoicesArgs.startDate, "start-date", "", "Data początkowa")
	SyncInvoicesCommand.FlagSet.StringVar(&syncInvoicesArgs.params.IssuerToken, "token", "", "Token sesji interaktywnej lub nazwa zmiennej środowiskowej która go zawiera")
	// SyncInvoicesCommand.FlagSet.StringVar(&syncInvoicesArgs.params.Token, "token", "", "Token sesji")
	SyncInvoicesCommand.FlagSet.StringVar(&syncInvoicesArgs.refresh, "refresh", "", "odświeża istniejący rejestr faktur według istniejącego pliku")

	registerCommand(&SyncInvoicesCommand.Command)
}

func syncInvoicesRun(c *Command) error {
	var err error

	// is it a refresh operation?
	if syncInvoicesArgs.refresh != "" {
		registry, err := registryPkg.LoadRegistry(syncInvoicesArgs.refresh)
		if err != nil {
			return fmt.Errorf("nie udało się załadować pliku rejestru: %v", err)
		}
		syncInvoicesArgs.params.DestPath = path.Dir(syncInvoicesArgs.refresh)

		apiClient, err := client.APIClient_Init(registry.Environment)
		if err != nil {
			return fmt.Errorf("nieznane środowisko: %v", registry.Environment)
		}

		return sync.SyncInvoices(apiClient, &syncInvoicesArgs.params, registry)
	}

	// is it a new request?
	if syncInvoicesArgs.params.DestPath == "" || syncInvoicesArgs.startDate == "" ||
		(!syncInvoicesArgs.params.Income &&
			!syncInvoicesArgs.params.Cost &&
			!syncInvoicesArgs.params.Subject3 &&
			!syncInvoicesArgs.params.SubjectAuthorized) {
		c.FlagSet.Usage()
		return nil
	}

	if syncInvoicesArgs.params.StartDate, err = time.ParseInLocation("2006-01-02", syncInvoicesArgs.startDate, time.Now().Location()); err != nil {
		return fmt.Errorf("invalid date supplied: %s", syncInvoicesArgs.startDate)
	}

	environment := client.ProductionEnvironment
	if syncInvoicesArgs.testGateway {
		environment = client.TestEnvironment
	}

	gateway, err := client.APIClient_Init(environment)
	if err != nil {
		return fmt.Errorf("nieznane środowisko: %v", environment)
	}

	return sync.SyncInvoices(gateway, &syncInvoicesArgs.params, nil)
}
