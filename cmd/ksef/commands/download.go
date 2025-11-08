package commands

import (
	"fmt"
	"ksef/cmd/ksef/flags"
	"ksef/internal/config"
	environmentPkg "ksef/internal/environment"
	registryPkg "ksef/internal/registry"
	v2 "ksef/internal/sei/api/client/v2"
	"ksef/internal/sei/api/client/v2/auth/validator"
	"ksef/internal/sei/api/client/v2/invoices"
	typesInvoices "ksef/internal/sei/api/client/v2/types/invoices"
	"path"
	"time"

	"github.com/spf13/cobra"
)

var syncInvoicesCommand = &cobra.Command{
	Use:   "download",
	Short: "Pobranie listy faktur",
	RunE:  syncInvoices,
}

type syncInvoicesArgsType struct {
	params        invoices.SyncParams
	startDate     time.Time
	startDateType registryPkg.DateType
	subjectNIP    string
	refresh       string
}

var syncInvoicesArgs = &syncInvoicesArgsType{
	startDateType: registryPkg.DateTypeStorage,
}

func init() {
	flagSet := syncInvoicesCommand.Flags()

	flagSet.BoolFunc("income", "pobranie faktur przychodowych (type=Subject1)", func(s string) error {
		syncInvoicesArgs.params.SubjectType = typesInvoices.SubjectTypeIssuer
		return nil
	})
	flagSet.BoolFunc("cost", "pobranie faktur kosztowych (type=Subject2)", func(s string) error {
		syncInvoicesArgs.params.SubjectType = typesInvoices.SubjectTypeRecipient
		return nil
	})
	flagSet.BoolFunc("payer", "pobranie faktur podmiotu innego (type=Subject3)", func(s string) error {
		syncInvoicesArgs.params.SubjectType = typesInvoices.SubjectTypePayer
		return nil
	})
	flagSet.BoolVarP(&syncInvoicesArgs.params.PDF, "pdf", "p", false, "generuj PDF dla pobranych faktur")
	flagSet.StringVarP(&syncInvoicesArgs.subjectNIP, "nip", "n", "", "numer NIP podmiotu")
	flagSet.IntVarP(&syncInvoicesArgs.params.PageSize, "page-size", "", 50, "liczba faktur na stronę odpowiedzi")
	flagSet.TimeVarP(&syncInvoicesArgs.startDate, "start-date", "s", time.Now().Truncate(24*time.Hour), []string{"2006-01-02"}, "data początkowa")
	flagSet.VarP(flags.StringChoice([]string{
		string(registryPkg.DateTypeIssue),
		string(registryPkg.DateTypeInvoicing),
		string(registryPkg.DateTypeStorage),
	}), "date-type", "", "typ daty używany do odpytywania listy faktur")
	flagSet.StringVarP(&syncInvoicesArgs.params.DestPath, "registry", "o", "", "katalog rejestru (zostanie stworzony jeśli nie istnieje)")
	flagSet.StringVarP(&syncInvoicesArgs.refresh, "refresh", "r", "", "odświeża istniejący rejestr faktur")

	syncInvoicesCommand.Flags().SortFlags = false
	syncInvoicesCommand.MarkFlagsMutuallyExclusive("income", "cost", "payer")
}

func syncInvoices(cmd *cobra.Command, _ []string) error {
	var err error
	var initializers []v2.InitializerFunc
	var registry *registryPkg.InvoiceRegistry
	var authValidator validator.AuthChallengeValidator
	var env = environmentPkg.FromContext(cmd.Context())

	// is it a refresh operation?
	if syncInvoicesArgs.refresh != "" {
		registry, err := registryPkg.LoadRegistry(syncInvoicesArgs.refresh)
		if err != nil {
			return fmt.Errorf("nie udało się załadować pliku rejestru: %v", err)
		}
		syncInvoicesArgs.params.DestPath = path.Dir(syncInvoicesArgs.refresh)
		env = registry.Environment
		if authValidator, err = authChallengeValidatorInstance(
			cmd, registry.Issuer, env,
		); err != nil {
			return err
		}
	} else {
		// is it a new request?
		if syncInvoicesArgs.subjectNIP == "" ||
			syncInvoicesArgs.params.DestPath == "" ||
			syncInvoicesArgs.startDate.IsZero() ||
			syncInvoicesArgs.params.SubjectType == typesInvoices.SubjectTypeInvalid {
			return cmd.Help()
		}

		registry, err = registryPkg.OpenOrCreate(syncInvoicesArgs.params.DestPath)
		if err != nil {
			return err
		}

		registry.QueryCriteria.SubjectType = string(syncInvoicesArgs.params.SubjectType)
		registry.QueryCriteria.DateFrom = syncInvoicesArgs.startDate
		registry.QueryCriteria.DateType = syncInvoicesArgs.startDateType
		registry.Issuer = syncInvoicesArgs.subjectNIP
		registry.Environment = env

		if authValidator, err = authChallengeValidatorInstance(
			cmd, syncInvoicesArgs.subjectNIP, env,
		); err != nil {
			return err
		}
	}

	initializers = append(initializers, v2.WithRegistry(registry), v2.WithAuthValidator(authValidator))

	cli, err := v2.NewClient(cmd.Context(), config.GetConfig(), env, initializers...)
	if err != nil {
		return fmt.Errorf("błąd inicjalizacji klienta: %v", err)
	}
	defer cli.Close()

	return cli.SyncInvoices(cmd.Context(), syncInvoicesArgs.params)

}
