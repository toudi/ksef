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

func AuthMethod(cmd *cobra.Command) {
	var flagSet = cmd.Flags()
	flagSet.Bool(FlagNameSessionToken, false, "postaraj się użyć tokenów sesyjnych")
	flagSet.String(FlagNameCertPath, "", "ścieżka do certyfikatu używanego do autoryzacji")
	flagSet.String(FlagNameKSeFToken, "", "token KSeF lub nazwa zmiennej srodowiskowej która go zawiera")

	cmd.MarkFlagsOneRequired(FlagNameCertPath, FlagNameKSeFToken, FlagNameSessionToken)
}

func Logout(cmd *cobra.Command) {
	cmd.Flags().Bool(FlagNameLogout, false, "wyloguj sesję po zakończeniu operacji")
}
