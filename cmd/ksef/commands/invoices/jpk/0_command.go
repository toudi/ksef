package jpk

import (
	"errors"
	"ksef/cmd/ksef/flags"
	"ksef/internal/invoicesdb/jpk"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var JPKCommand = &cobra.Command{
	Use:     "jpk",
	Short:   "wygeneruj dokument JPK na podstawie bazy faktur",
	Args:    cobra.MaximumNArgs(2),
	RunE:    generateJPK,
	PreRunE: checkMonthArgs,
}

func checkMonthArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 && len(args) < 2 {
		return errors.New("użycie: jpk rok miesiac")
	}
	return nil
}

func init() {
	flagSet := JPKCommand.Flags()
	flags.NIP(flagSet)
	JPKCommand.MarkFlagRequired(flags.FlagNameNIP)
	JPKCommand.AddCommand(jpkExclude, jpk50PercVAT, jpkFixedAssets, jpkClearFlags, jpkDumpItems)
}

func generateJPK(cmd *cobra.Command, args []string) error {
	vip := viper.GetViper()
	var month time.Time = time.Now().Local().AddDate(0, -1, 0)
	if len(args) > 0 {
		yearNum, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		monthNum, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		month = time.Date(yearNum, time.Month(monthNum), 1, 0, 0, 0, 0, time.Local)
	}

	jpk, err := jpk.NewJPK(month, vip)
	if err != nil {
		return err
	}
	return jpk.Generate()
}
