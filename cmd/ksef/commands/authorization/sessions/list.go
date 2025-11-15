package sessions

import (
	"fmt"
	"ksef/cmd/ksef/flags"
	"ksef/internal/client/v2/auth"
	"ksef/internal/runtime"

	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var GetAuthSessions = &cobra.Command{
	Use:               "sessions",
	Short:             "wyświetla listę sesji dla aktualnego tokenu sesyjnego",
	PersistentPreRunE: getSessions,
	RunE:              runGetAuthSessions,
}

func init() {
	flags.NIP(GetAuthSessions.PersistentFlags())
	GetAuthSessions.AddCommand(closeSession)
}

var activeSessions *auth.AuthSessionsResponse
var tokenManager *auth.TokenManager

func getSessions(cmd *cobra.Command, _ []string) error {
	var err error

	vip := viper.GetViper()
	nip, err := cmd.Flags().GetString(flags.FlagNameNIP)
	if err != nil {
		return err
	}
	runtime.SetNIP(vip, nip)

	tokenManager, err = auth.NewTokenManager(
		cmd.Context(),
		vip,
		nil,
	)
	if err != nil {
		return err
	}

	activeSessions, err = tokenManager.GetAuthSessions(cmd.Context())
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
