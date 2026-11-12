package sessions

import (
	"fmt"

	"ksef/cmd/ksef/flags"
	"ksef/internal/resources"
	"ksef/internal/runtime"

	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var GetAuthSessions = &cobra.Command{
	Use:   "sessions",
	Short: "wyświetla listę sesji dla aktualnego tokenu sesyjnego",
	RunE:  runGetAuthSessions,
}

func init() {
	flags.NIP(GetAuthSessions.PersistentFlags())
	GetAuthSessions.AddCommand(closeSession)
}

func runGetAuthSessions(cmd *cobra.Command, _ []string) (commandErr error) {
	ctx := cmd.Context()
	vip := viper.GetViper()

	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return err
	}

	ksefClient, err := resources.Client(ctx, vip, nip)
	if err != nil {
		return err
	}
	defer func() {
		resources.Close(ctx, commandErr)
	}()

	sessions, err := ksefClient.AuthSessions(ctx)
	if err != nil {
		return err
	}

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "id"},
			{Align: simpletable.AlignCenter, Text: "aktualna"},
			{Align: simpletable.AlignCenter, Text: "rozpoczęta"},
			{Align: simpletable.AlignCenter, Text: "refresh token ważny do"},
		},
	}

	for _, session := range sessions.Sessions {
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
