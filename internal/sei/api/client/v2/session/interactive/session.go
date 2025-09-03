package interactive

import (
	"context"
	"errors"
	"ksef/internal/config"
	HTTP "ksef/internal/http"
	"ksef/internal/logging"
	"ksef/internal/registry"
	"ksef/internal/sei/api/client/v2/auth"
	"time"
)

type Session struct {
	collection          *registry.InvoiceCollection
	finished            bool
	tokenUpdatesChannel chan auth.TokenUpdate
	sessionToken        string
	httpClient          HTTP.Client
	initialized         bool
	apiConfig           config.APIConfig
}

var ErrObtainSessionTokenTimeout = errors.New("timeout waiting for session token")

func NewSession(httpClient HTTP.Client, tokenUpdatesChannel chan auth.TokenUpdate, collection *registry.InvoiceCollection) *Session {
	return &Session{
		tokenUpdatesChannel: tokenUpdatesChannel,
		httpClient:          httpClient,
		collection:          collection,
	}
}

func (s *Session) UploadInvoices() error {
	var initializedC = make(chan bool)

	go s.eventLoop(initializedC)
	defer func() {
		s.finished = true
	}()

	// we have to wait for the session token to be ready but let's not wait indefinetely
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	select {
	case <-ctx.Done():
		cancel()
		return ErrObtainSessionTokenTimeout
	case <-initializedC:
		logging.AuthLogger.Debug("session token obtained")
		cancel()
	}

	// v2 specs forces us to group invoices by their form code
	// on the other hand, it no longer forces us to send invoices through a 3rd party server
	for formCode, files := range s.collection.Files {
		if err := s.uploadInvoicesForForm(context.Background(), formCode, files); err != nil {
			return err
		}
	}

	return nil
}

// eventLoop manages the session's lifecycle. it's required because the session token validation is
// asynchronous and so it is the token manager that needs to notify the session that the token is
// indeed ready (and what it is). This event will be called regardless whether the token was
// established for the first time or it was refreshed
func (s *Session) eventLoop(initialized chan bool) {
	var heartbeat = time.NewTicker(time.Second)

	for !s.finished {
		select {
		case <-heartbeat.C:
			// nothing to do here, really
			if !s.initialized {
				logging.InteractiveLogger.Debug("interactive session not initialized, waiting for the token")
			}
		case tokenUpdate := <-s.tokenUpdatesChannel:
			if tokenUpdate.Err != nil {
				logging.InteractiveLogger.Error("error refreshing token", "error", tokenUpdate.Err)
				s.finished = true
				break
			}
			logging.InteractiveLogger.Debug("token refreshed")
			s.sessionToken = tokenUpdate.Token
			if !s.initialized {
				initialized <- true
			}
		}
	}

	heartbeat.Stop()
}

func (s *Session) authHeaders() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + s.sessionToken,
	}
}
