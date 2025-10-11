package authorization

import (
	"ksef/cmd/ksef/commands/authorization/xades"

	"github.com/spf13/cobra"
)

var AuthCommand = &cobra.Command{
	Use:   "auth",
	Short: "zarządzanie autoryzacją KSeF",
}

var nip string

func init() {
	// flags.NIPVarP(AuthCommand.PersistentFlags())
	AuthCommand.PersistentFlags().StringVarP(&nip, "nip", "n", "", "numer NIP podmiotu")

	AuthCommand.AddCommand(xades.XadesCommand)
}
