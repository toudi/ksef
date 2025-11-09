package download

import (
	"ksef/cmd/ksef/flags"
	"ksef/internal/client/v2/types/invoices"
	"time"

	"github.com/spf13/pflag"
)

const (
	flagIncremental = "incremental"
	flagPDF         = "pdf"
	flagStartDate   = "start-date"
	flagEndDate     = "end-date"
	flagPageSize    = "page-size"
	flagDateType    = "date-type"
)

var subjectTypes []invoices.SubjectType

func registerFlags(flagSet *pflag.FlagSet) {
	var today = time.Now().Local()
	var firstDayOfMonth = today.AddDate(0, 0, -today.Day()+1)

	flags.NIP(flagSet)
	flagSet.BoolFunc("income", "pobieranie faktur przychodowych (Subject=Subject1)", func(s string) error {
		subjectTypes = append(subjectTypes, invoices.SubjectTypeIssuer)
		return nil
	})
	flagSet.BoolFunc("cost", "pobieranie faktur kosztowych (Subject=Subject2)", func(s string) error {
		subjectTypes = append(subjectTypes, invoices.SubjectTypeRecipient)
		return nil
	})
	flagSet.BoolFunc("payer", "pobieranie faktur płatnika (Subject=Subject3)", func(s string) error {
		subjectTypes = append(subjectTypes, invoices.SubjectTypePayer)
		return nil
	})
	flagSet.BoolFunc("authorized", "pobieranie faktur strony upoważnionej (Subject=SubjectAuthorized)", func(s string) error {
		subjectTypes = append(subjectTypes, invoices.SubjectTypeAuthorized)
		return nil
	})
	flagSet.Bool(flagPDF, false, "generuj PDF dla pobranych faktur")
	flagSet.BoolP(flagIncremental, "i", false, "pobieranie przyrostowe")
	flagSet.String(flagStartDate, firstDayOfMonth.Format("2006-01-02"), "data początkowa")
	flagSet.String(flagEndDate, "", "data końcowa")
	flagSet.IntP(flagPageSize, "", 50, "liczba faktur na stronę odpowiedzi")
	flagSet.Var(flags.StringChoice([]string{
		string(invoices.DateTypeIssue),
		string(invoices.DateTypeInvoicing),
		string(invoices.DateTypeStorage),
	}), flagDateType, "typ daty używany do odpytywania listy faktur")

	flagSet.SortFlags = false
}

func getDownloadParams(flagSet *pflag.FlagSet) (params invoices.DownloadParams, err error) {
	if params.Incremental, err = flagSet.GetBool(flagIncremental); err != nil {
		return params, err
	}
	if params.PDF, err = flagSet.GetBool(flagPDF); err != nil {
		return params, err
	}
	params.SubjectTypes = subjectTypes
	startDate, err := flagSet.GetString(flagStartDate)
	if err != nil {
		return params, err
	}
	if params.StartDate, err = time.ParseInLocation(time.DateOnly, startDate, time.Local); err != nil {
		return params, err
	}
	endDate, err := flagSet.GetString(flagEndDate)
	if err != nil {
		return params, err
	}
	if endDate != "" {
		if *params.EndDate, err = time.ParseInLocation(time.DateOnly, endDate, time.Local); err != nil {
			return params, err
		}
	}

	return params, nil
}
