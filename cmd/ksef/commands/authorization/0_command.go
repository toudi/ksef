package authorization

import (
	"ksef/cmd/ksef/commands/authorization/sessions"
	"ksef/internal/runtime"

	"github.com/spf13/cobra"
)

var AuthCommand = &cobra.Command{
	Use:   "auth",
	Short: "zarządzanie autoryzacją KSeF",
}

func init() {
	runtime.CertProfileFlag(AuthCommand.PersistentFlags())
	AuthCommand.AddCommand(sessions.GetAuthSessions)
}
