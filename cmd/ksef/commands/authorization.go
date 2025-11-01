package commands

import (
	"flag"
)

const (
	authBackendKsefToken      = "token"
	authBackendXadesSignature = "xades:epuap"
	authBackendXadesKsefCert  = "xades:ksef-cert"
)

// common module for setting up authentication backend properties
type authParamsKsefToken struct {
	issuerToken string
}

type authParamsXadesSignature struct {
	signedRequestFile string
}

type authParamsT struct {
	backend              string
	ksefTokenParams      authParamsKsefToken
	xadesSignatureParams authParamsXadesSignature
}

var (
	authParams = authParamsT{
		backend: authBackendKsefToken,
	}
	// environment environmentPkg.Environment = environmentPkg.Production
)

func initAuthParams(flagSet *flag.FlagSet) {
	flagSet.StringVar(&authParams.backend, "b", authBackendKsefToken, "Backend autoryzacji")
	flagSet.StringVar(&authParams.ksefTokenParams.issuerToken, "ksef:token", "", "Token sesji interaktywnej lub nazwa zmiennej środowiskowej która go zawiera (tylko jeśli wybranym backendem jest autoryzacja tokenem KSeF)")
	flagSet.StringVar(&authParams.xadesSignatureParams.signedRequestFile, "xades:signed-challenge", "", "**PODPISANY** plik wyzwania XML. Aby wygenerować wyzwanie, użyj komendy xades-init")
}

func testGatewayFlag(flagSet *flag.FlagSet) {
	flagSet.BoolFunc("t", "Użyj bramki testowej", func(s string) error {
		// environment = environmentPkg.Test
		return nil
	})
}

// func authValidatorInstance(issuerNip string) validator.AuthChallengeValidator {
// 	// harcode challenge validator to ksef token for now
// 	// TODO: implement XADeS
// 	tokenValidator := kseftoken.NewKsefTokenHandler(
// 		config.GetConfig().APIConfig(environment), issuerNip,
// 	)
// 	if authParams.ksefTokenParams.issuerToken != "" {
// 		logging.AuthLogger.Warn("overriding KSeF token")
// 		// TODO: remember to properly handle casting to interface error once XaDes auth validator is implemented
// 		tokenValidatorInstance := tokenValidator.(*kseftoken.KsefTokenHandler)
// 		tokenValidatorInstance.SetKsefToken(authParams.ksefTokenParams.issuerToken)
// 	}
// 	return tokenValidator
// }
