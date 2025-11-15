package flags

import (
	"github.com/spf13/cobra"
)

const (
	FlagNameCertPath        = "auth.cert"
	FlagNameSessionToken    = "auth.token"
	FlagNameKSeFToken       = "auth.ksef-token"
	FlagNameSignedChallenge = "signed"
	FlagNameLogout          = "logout"
)

func SignedChallenge(cmd *cobra.Command) {
	cmd.Flags().StringP(FlagNameSignedChallenge, "s", "", "lokalizacja *PODPISANEGO* pliku wyzwania")
}

func Logout(cmd *cobra.Command) {
	cmd.Flags().Bool(FlagNameLogout, false, "wyloguj sesję po zakończeniu operacji")
}
