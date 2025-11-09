package token

import (
	"bytes"
	"context"
	"io"
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/auth/validator"
	"ksef/internal/client/v2/auth/xades"
	"ksef/internal/config"
	"ksef/internal/http"
	"ksef/internal/logging"
	"os"

	"github.com/zalando/go-keyring"
)

type mode uint8

const (
	modeRegular       mode = 0x01
	modeDumpChallenge mode = 0x02
	modeUseSignedFile mode = 0x03
)

type TokenHandler struct {
	gateway             config.Gateway
	httpClient          *http.Client
	nip                 string
	eventChannel        chan validator.AuthEvent
	certsDB             *certsdb.CertificatesDB
	mode                mode
	challengeDumpPath   string
	signedChallengeFile string
}

func NewAuthHandler(gateway config.Gateway, nip string, initializers ...initializerFunc) validator.AuthChallengeValidator {
	handler := &TokenHandler{
		gateway:      gateway,
		eventChannel: make(chan validator.AuthEvent),
		nip:          nip,
		mode:         modeRegular,
	}

	for _, initializer := range initializers {
		initializer(handler)
	}

	return handler
}

func (e *TokenHandler) Event() chan validator.AuthEvent {
	return e.eventChannel
}

func (e *TokenHandler) sendAuthEvent(event validator.AuthEvent) {
	e.eventChannel <- event
}

func (e *TokenHandler) initialize(ctx context.Context, errChan chan error) {
	if e.mode == modeDumpChallenge {
		go e.sendAuthEvent(validator.AuthEvent{
			State: validator.StateInitialized,
		})
		errChan <- nil
		return
	}
	if e.mode == modeRegular {
		sessionTokens, err := keyring.Get(string(e.gateway)+"-sessionTokens", e.nip)
		if err == nil {
			if sessionTokens != "" {
				go e.sendAuthEvent(validator.AuthEvent{
					State:         validator.StateTokensRestored,
					SessionTokens: sessionTokens,
				})
				errChan <- nil
				return
			}
			// if token retrieval would have been successful the function would
			// have already returned. let's proceed with regular work - obtaining
			// challenge
			go e.sendAuthEvent(validator.AuthEvent{
				State: validator.StateInitialized,
			})
		} else {
			logging.AuthLogger.Warn("nie udało się odczytać tokenów sesyjnych")
		}
		errChan <- nil
		return
	}
	if e.mode == modeUseSignedFile {
		signedChallenge, err := os.Open(e.signedChallengeFile)
		if err != nil {
			errChan <- nil
			return
		}
		defer signedChallenge.Close()

		errChan <- validateSignedChallenge(
			ctx,
			e.httpClient,
			signedChallenge,
			func(resp validator.ValidationReference) {
				e.eventChannel <- validator.AuthEvent{
					State:               validator.StateValidationReferenceResult,
					ValidationReference: &resp,
				}
			},
		)
		return
	}
}

func (e *TokenHandler) Initialize(ctx context.Context, httpClient *http.Client) error {
	e.httpClient = httpClient
	var errChan = make(chan error)
	go e.initialize(ctx, errChan)
	err := <-errChan
	return err
}

func (e *TokenHandler) ValidateChallenge(ctx context.Context, challenge validator.AuthChallenge) error {
	// we have our challenge. we now need to sign it and send using validateSignedChallenge
	var sourceDocument io.ReadWriter = new(bytes.Buffer)
	var err error

	if e.mode == modeDumpChallenge {
		logging.AuthLogger.Debug("żądanie zrzucenia wyzwania do pliku", "output", e.challengeDumpPath)
		if sourceDocument, err = os.Create(e.challengeDumpPath); err != nil {
			logging.AuthLogger.Error("nie udało się utworzyć pliku", "err", err)
			return err
		}
		defer sourceDocument.(*os.File).Close()
	}

	if err = dumpChallengeToWriter(challenge, e.nip, sourceDocument); err != nil {
		logging.AuthLogger.Error("nie udało się zapisać wyzwania do bufora", "err", err)
		return err
	}

	if e.mode == modeDumpChallenge {
		// in this mode we only want to dump challenge to file and exit
		go e.sendAuthEvent(validator.AuthEvent{
			State: validator.StateExit,
		})
		return nil
	}
	// great. now we can sign it using the certificate
	certificate, err := e.certsDB.GetByUsage(certsdb.UsageAuthentication, e.nip)
	if err != nil {
		return err
	}
	var signedDocument bytes.Buffer
	if err = xades.SignAuthChallenge(sourceDocument, certificate, &signedDocument); err != nil {
		return err
	}
	// perfect. final step - let's post it to the validation endpoint
	return validateSignedChallenge(
		ctx,
		e.httpClient,
		bytes.NewReader(signedDocument.Bytes()),
		func(resp validator.ValidationReference) {
			e.eventChannel <- validator.AuthEvent{
				State:               validator.StateValidationReferenceResult,
				ValidationReference: &resp,
			}
		},
	)
}
