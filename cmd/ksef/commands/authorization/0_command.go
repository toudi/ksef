package authorization

import (
	"ksef/cmd/ksef/commands/authorization/xades"
	"ksef/cmd/ksef/flags"

	"github.com/spf13/cobra"
)

var AuthCommand = &cobra.Command{
	Use:   "auth",
	Short: "zarządzanie autoryzacją KSeF",
}

func init() {
	flags.NIP(AuthCommand.PersistentFlags())

	AuthCommand.AddCommand(xades.XadesCommand)
}
