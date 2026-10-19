package sessions

import (
	"fmt"
	"ksef/cmd/ksef/flags"
	"ksef/internal/client/v2/auth"
	"ksef/internal/logging"
	"ksef/internal/resources"

	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var GetAuthSessions = &cobra.Command{
	Use:     "sessions",
	Short:   "wyświetla listę sesji dla aktualnego tokenu sesyjnego",
	PreRunE: getSessions,
	RunE:    runGetAuthSessions,
}

func init() {
	flags.NIP(GetAuthSessions.PersistentFlags())
	GetAuthSessions.AddCommand(closeSession)
}

var (
	activeSessions *auth.AuthSessionsResponse
	tokenManager   *auth.TokenManager
)

func getSessions(cmd *cobra.Command, _ []string) error {
	var err error

	vip := viper.GetViper()
	ctx := cmd.Context()

	keyring, err := resources.Keyring(vip)
	if err != nil {
		logging.SeiLogger.Error("błąd inicjalizacji keyringu", "err", err)
		return err
	}
	certsDB, err := resources.CertsDB(vip)
	if err != nil {
		return err
	}
	defer resources.Close(ctx)

	tokenManager, err = auth.NewTokenManager(
		cmd.Context(),
		vip,
		nil,
		keyring,
		certsDB,
	)
	if err != nil {
		return err
	}

	activeSessions, err = tokenManager.GetAuthSessions(ctx)
	if err != nil {
		return err
	}

	return nil
}

func runGetAuthSessions(cmd *cobra.Command, _ []string) error {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "id"},
			{Align: simpletable.AlignCenter, Text: "aktualna"},
			{Align: simpletable.AlignCenter, Text: "rozpoczęta"},
			{Align: simpletable.AlignCenter, Text: "refresh token ważny do"},
		},
	}

	for _, session := range activeSessions.Sessions {
		table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
			{
				Text: session.ID,
			},
			{
				Text: fmt.Sprintf("%t", session.Current),
			},
			{
				Text: session.StartDate.String(),
			},
			{
				Text: session.RefreshTokenValidUntil.String(),
			},
		})
	}

	table.SetStyle(simpletable.StyleDefault)
	table.Println()

	return nil
}
