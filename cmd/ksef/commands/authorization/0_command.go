package authorization

import (
	"ksef/cmd/ksef/commands/authorization/sessions"

	"github.com/spf13/cobra"
)

var AuthCommand = &cobra.Command{
	Use:   "auth",
	Short: "zarządzanie autoryzacją KSeF",
}

func init() {
	AuthCommand.AddCommand(sessions.GetAuthSessions)
}
