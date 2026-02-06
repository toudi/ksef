package config

import (
	"ksef/cmd/ksef/flags"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/utils"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

func prefixedFlag(prefix string, flagName string) string {
	if prefix != "" {
		return prefix + "." + flagName
	}
	return flagName
}

func DownloaderFlags(flagSet *pflag.FlagSet, prefix string) {
	today := time.Now().Local()
	firstDayOfMonth := today.AddDate(0, 0, -today.Day()+1)

	flagSet.BoolFunc(prefixedFlag(prefix, "income"), "pobieranie faktur przychodowych (Subject=Subject1)", func(s string) error {
		subjectTypes = append(subjectTypes, invoices.SubjectTypeIssuer)
		return nil
	})
	flagSet.BoolFunc(prefixedFlag(prefix, "cost"), "pobieranie faktur kosztowych (Subject=Subject2)", func(s string) error {
		subjectTypes = append(subjectTypes, invoices.SubjectTypeRecipient)
		return nil
	})
	flagSet.BoolFunc(prefixedFlag(prefix, "payer"), "pobieranie faktur płatnika (Subject=Subject3)", func(s string) error {
		subjectTypes = append(subjectTypes, invoices.SubjectTypePayer)
		return nil
	})
	flagSet.BoolFunc(prefixedFlag(prefix, "authorized"), "pobieranie faktur strony upoważnionej (Subject=SubjectAuthorized)", func(s string) error {
		subjectTypes = append(subjectTypes, invoices.SubjectTypeAuthorized)
		return nil
	})
	flagSet.Bool(prefixedFlag(prefix, flagPDF), false, "generuj PDF dla pobranych faktur")
	flagSet.BoolP(prefixedFlag(prefix, flagIncremental), "i", false, "pobieranie przyrostowe")
	flagSet.String(prefixedFlag(prefix, flagStartDate), firstDayOfMonth.Format("2006-01-02"), "data początkowa")
	flagSet.String(prefixedFlag(prefix, flagEndDate), "", "data końcowa")
	flagSet.IntP(prefixedFlag(prefix, flagPageSize), "", 50, "liczba faktur na stronę odpowiedzi")
	flagSet.Var(flags.StringChoice([]string{
		string(invoices.DateTypeIssue),
		string(invoices.DateTypeInvoicing),
		string(invoices.DateTypeStorage),
	}), prefixedFlag(prefix, flagDateType), "typ daty używany do odpytywania listy faktur")

	flagSet.SortFlags = false
}

func GetDownloaderConfig(vip *viper.Viper, prefix string) (params invoices.DownloadParams, err error) {
	params.Incremental = vip.GetBool(prefixedFlag(prefix, flagIncremental))
	params.PDF = vip.GetBool(prefixedFlag(prefix, flagPDF))
	params.SubjectTypes = subjectTypes
	params.PageSize = vip.GetInt(prefixedFlag(prefix, flagPageSize))
	dateRangeType := vip.GetString(prefixedFlag(prefix, flagDateType))
	params.DateType = invoices.DateRangeType(dateRangeType)
	if params.Incremental || params.DateType == "" {
		params.DateType = invoices.DateTypeStorage
	}
	startDate, err := utils.ParseTimeFromString(vip.GetString(prefixedFlag(prefix, flagStartDate)))
	if err != nil {
		return params, err
	}
	params.StartDate = startDate
	endDate := vip.GetString(prefixedFlag(prefix, flagEndDate))
	if endDate != "" {
		parsedEndDate, err := utils.ParseTimeFromString(endDate)
		if err != nil {
			return params, err
		}
		params.EndDate = &parsedEndDate
	}

	return params, nil
}
