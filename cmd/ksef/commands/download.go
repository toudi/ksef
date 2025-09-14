package commands

import (
	"flag"
	"fmt"
	"ksef/internal/config"
	registryPkg "ksef/internal/registry"
	v2 "ksef/internal/sei/api/client/v2"
	"ksef/internal/sei/api/client/v2/auth/validator"
	"ksef/internal/sei/api/client/v2/invoices"
	"path"
	"time"
)

type syncInvoicesCommand struct {
	Command
}

type syncInvoicesArgsType struct {
	params         invoices.SyncParams
	startDateInput string
	startDate      time.Time
	subjectNIP     string
	refresh        string
}

var (
	SyncInvoicesCommand *syncInvoicesCommand
	syncInvoicesArgs    = &syncInvoicesArgsType{}
)

func init() {
	SyncInvoicesCommand = &syncInvoicesCommand{
		Command: Command{
			Name:        "download",
			FlagSet:     flag.NewFlagSet("download", flag.ExitOnError),
			Description: "Synchronizuje listę faktur z KSeF do katalogu lokalnego",
			Run:         syncInvoicesRun,
		},
	}

	flagSet := SyncInvoicesCommand.FlagSet
	initAuthParams(flagSet)
	testGatewayFlag(flagSet)
	flagSet.StringVar(&syncInvoicesArgs.params.DestPath, "d", "", "Katalog docelowy")

	// specify subject type
	flagSet.BoolFunc("income", "Synchronizuj faktury przychodowe (Podmiot1)", func(s string) error {
		syncInvoicesArgs.params.SubjectType = invoices.SubjectTypeIssuer
		return nil
	})
	flagSet.BoolFunc("cost", "Synchronizuj faktury kosztowe (Podmiot2)", func(s string) error {
		syncInvoicesArgs.params.SubjectType = invoices.SubjectTypeRecipient
		return nil
	})
	flagSet.BoolFunc("payer", "Synchronizuj faktury podmiotu innego (Podmiot3)", func(s string) error {
		syncInvoicesArgs.params.SubjectType = invoices.SubjectTypePayer
		return nil
	})
	flagSet.BoolFunc("subjectAuthorized", "Synchronizuj faktury podmiotu upoważnionego (???)", func(s string) error {
		syncInvoicesArgs.params.SubjectType = invoices.SubjectTypeAuthorized
		return nil
	})

	flagSet.BoolVar(&syncInvoicesArgs.params.PDF, "pdf", false, "Generuj PDF dla pobranych faktur")
	flagSet.IntVar(&syncInvoicesArgs.params.PageSize, "page-size", 50, "Liczba faktur na stronę odpowiedzi")
	flagSet.StringVar(&syncInvoicesArgs.subjectNIP, "nip", "", "Numer NIP podmiotu")
	flagSet.StringVar(&syncInvoicesArgs.startDateInput, "start-date", "", "Data początkowa")
	flagSet.StringVar(&syncInvoicesArgs.refresh, "refresh", "", "odświeża istniejący rejestr faktur według istniejącego pliku")

	registerCommand(&SyncInvoicesCommand.Command)
}

func syncInvoicesRun(c *Command) error {
	var err error
	var initializers []v2.InitializerFunc
	var registry *registryPkg.InvoiceRegistry
	var authValidator validator.AuthChallengeValidator

	// is it a refresh operation?
	if syncInvoicesArgs.refresh != "" {
		registry, err := registryPkg.LoadRegistry(syncInvoicesArgs.refresh)
		if err != nil {
			return fmt.Errorf("nie udało się załadować pliku rejestru: %v", err)
		}
		syncInvoicesArgs.params.DestPath = path.Dir(syncInvoicesArgs.refresh)
		environment = registry.Environment
		authValidator = authValidatorInstance(registry.Issuer)
	} else {
		// is it a new request?
		if syncInvoicesArgs.subjectNIP == "" ||
			syncInvoicesArgs.params.DestPath == "" ||
			syncInvoicesArgs.startDateInput == "" ||
			syncInvoicesArgs.params.SubjectType == invoices.SubjectTypeInvalid {
			c.FlagSet.Usage()
			return nil
		}

		if syncInvoicesArgs.startDate, err = time.ParseInLocation("2006-01-02", syncInvoicesArgs.startDateInput, time.Now().Location()); err != nil {
			return fmt.Errorf("invalid date supplied: %s", syncInvoicesArgs.startDate)
		}

		registry, err = registryPkg.OpenOrCreate(syncInvoicesArgs.params.DestPath)
		if err != nil {
			return err
		}

		registry.QueryCriteria.DateFrom = syncInvoicesArgs.startDate
		registry.Issuer = syncInvoicesArgs.subjectNIP

		authValidator = authValidatorInstance(syncInvoicesArgs.subjectNIP)
	}

	initializers = append(initializers, v2.WithRegistry(registry), v2.WithAuthValidator(authValidator))

	cli, err := v2.NewClient(c.Context, config.GetConfig(), environment, initializers...)
	if err != nil {
		return fmt.Errorf("błąd inicjalizacji klienta: %v", err)
	}

	return cli.SyncInvoices(c.Context, syncInvoicesArgs.params)
}
