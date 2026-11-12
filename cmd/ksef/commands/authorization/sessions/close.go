package sessions

import (
	"ksef/internal/logging"
	"ksef/internal/resources"
	"ksef/internal/runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	allSessions        = ":all:"
	nonCurrentSessions = ":nocurrent:"
)

var closeSession = &cobra.Command{
	Use:   "close",
	Short: "kończy wybraną sesję",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runCloseSession,
}

func runCloseSession(cmd *cobra.Command, ids []string) (commandErr error) {
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

	sessionToClose := nonCurrentSessions
	if len(ids) > 0 {
		sessionToClose = ids[0]
	}

	for _, session := range sessions.Sessions {
		// if we've specified a single session ID this if statement won't be touched
		if session.ID != sessionToClose {
			// this can mean either :all: or :nocurrent:
			if session.Current && sessionToClose != allSessions {
				continue
			}
		}

		if err := ksefClient.LogoutSession(session.ID); err != nil {
			logging.AuthLogger.Error("błąd kończenia sesji", "err", err)
		}
	}

	return nil
}
