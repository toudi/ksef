package authorization

import (
	"errors"
	"ksef/cmd/ksef/flags"
	v2 "ksef/internal/client/v2"
	"ksef/internal/config"
	"ksef/internal/environment"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type bindNipArgsType struct {
	pesel string
}

var (
	bindNipArgs            bindNipArgsType
	errTestModeNotSelected = errors.New("komenda dostępna tylko przy użyciu testowej bramki KSeF")
	errNipIsRequired       = errors.New("numer NIP jest wymagany")
)

var bindNipCommand = &cobra.Command{
	Use:   "bind-nip",
	Short: "powiązanie numeru NIP z numerem PESEL w KSeF (działa tylko w trybie testowym)",
	RunE:  bindNipRun,
}

func init() {
	bindNipCommand.Flags().StringVarP(&bindNipArgs.pesel, "pesel", "p", "", "numer PESEL osoby upoważnionej")
	flags.NIP(bindNipCommand.Flags())

	AuthCommand.AddCommand(bindNipCommand)
}

func bindNipRun(cmd *cobra.Command, _ []string) error {
	env := environment.FromContext(cmd.Context())
	if env != environment.Test {
		return errTestModeNotSelected
	}
	nip, err := cmd.Flags().GetString("nip")
	if err != nil || nip == "" {
		return errNipIsRequired
	}

	cli, err := v2.NewClient(cmd.Context(), config.GetGateway(viper.GetViper()))
	if err != nil {
		return err
	}

	return cli.BindNIPToPESEL(cmd.Context(), nip, bindNipArgs.pesel)
}
