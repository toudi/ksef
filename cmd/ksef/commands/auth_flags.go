package commands

import (
	"errors"
	"ksef/cmd/ksef/flags"
	"ksef/internal/certsdb"
	"ksef/internal/config"
	environmentPkg "ksef/internal/environment"
	"ksef/internal/logging"
	kseftoken "ksef/internal/sei/api/client/v2/auth/ksef_token"
	"ksef/internal/sei/api/client/v2/auth/validator"
	"ksef/internal/sei/api/client/v2/auth/xades"

	"github.com/spf13/cobra"
)

var (
	errUnableToSelectAuthBackend = errors.New("Nie udało się wybrać metody autoryzacji")
)

func authChallengeValidatorInstance(cmd *cobra.Command, nip string, env environmentPkg.Environment) (validator.AuthChallengeValidator, error) {
	var certPath string
	var ksefToken string
	var err error
	var flagSet = cmd.Flags()

	if ksefToken, err = flagSet.GetString(flags.FlagNameKSeFToken); err != nil {
		return nil, err
	}
	if certPath, err = flagSet.GetString(flags.FlagNameCertPath); err != nil {
		return nil, err
	}

	apiConfig := config.GetConfig().APIConfig(env)
	certsDB, err := certsdb.OpenOrCreate(env)
	if err != nil {
		return nil, err
	}

	if certPath != "" {
		// cert-based authentication
		logging.AuthLogger.Debug("wybrano autoryzację certyfikatem kwalifikowanym")
		return xades.NewAuthHandler(
			apiConfig,
			nip,
			certPath,
		), nil
	} else if ksefToken != "" {
		logging.AuthLogger.Debug("wybrano autoryzację tokenem KSeF")
		// token-based authentication
		return kseftoken.NewKsefTokenHandler(
			apiConfig,
			certsDB,
			nip,
		), nil
	}

	return nil, errUnableToSelectAuthBackend
}
