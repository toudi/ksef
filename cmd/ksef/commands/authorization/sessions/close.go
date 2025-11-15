package sessions

import (
	"ksef/internal/logging"

	"github.com/spf13/cobra"
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

func runCloseSession(cmd *cobra.Command, ids []string) error {
	var sessionToClose = nonCurrentSessions
	if len(ids) > 0 {
		sessionToClose = ids[0]
	}

	for _, session := range activeSessions.Sessions {
		// if we've specified a single session ID this if statement won't be touched
		if session.ID != sessionToClose {
			// this can mean either :all: or :nocurrent:
			if session.Current && sessionToClose != allSessions {
				continue
			}
		}

		if err := tokenManager.Logout(session.ID); err != nil {
			logging.AuthLogger.Error("błąd kończenia sesji", "err", err)
		}
	}

	return nil
}
