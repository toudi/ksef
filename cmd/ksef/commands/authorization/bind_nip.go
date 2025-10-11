package authorization

import (
	"errors"
	"ksef/internal/environment"

	"github.com/spf13/cobra"
)

type bindNipArgsType struct {
	pesel string
}

var (
	bindNipArgs            bindNipArgsType
	errTestModeNotSelected = errors.New("komenda dostępna tylko przy użyciu testowej bramki KSeF")
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
	return nil
}
