package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTokenManager(t *testing.T) {
	t.Run("timeout reading token", func(t *testing.T) {
		t.Parallel()

		tm := &TokenManager{}

		token, err := tm.GetAuthorizationToken(1 * time.Second)
		require.ErrorIs(t, err, ErrTimeoutReadingToken)
		require.Empty(t, token)
	})

	t.Run("retrieve current token", func(t *testing.T) {
		t.Parallel()

		tm := &TokenManager{
			sessionTokens: &SessionTokens{
				AuthorizationToken: &TokenInfo{
					Token: "current-session-token",
				},
			},
		}

		token, err := tm.GetAuthorizationToken(1 * time.Second)
		require.NoError(t, err)
		require.Equal(t, "current-session-token", token)
	})
}
