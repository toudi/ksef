package authorization

import (
	"errors"
	"ksef/internal/config"
	"ksef/internal/environment"
	v2 "ksef/internal/sei/api/client/v2"

	"github.com/spf13/cobra"
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
	cfg := config.GetConfig()

	cli, err := v2.NewClient(cmd.Context(), cfg, env)
	if err != nil {
		return err
	}

	return cli.BindNIPToPESEL(cmd.Context(), nip, bindNipArgs.pesel)
}
