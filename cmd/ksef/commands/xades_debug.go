package commands

import (
	"flag"
	"ksef/internal/config"
	v2 "ksef/internal/sei/api/client/v2"
	"ksef/internal/sei/api/client/v2/auth/esignature"
)

type xadesDebugCommand struct {
	Command
}

type xadesDebugArgsType struct {
	signedFile string
}

var XadesDebugCommand *xadesDebugCommand
var xadesDebugArgs xadesDebugArgsType

func init() {
	XadesDebugCommand = &xadesDebugCommand{
		Command: Command{
			Name:        "xades-debug",
			Description: "inicjalizuje sesję przy pomocy podpisanego pliku autoryzacji a następnie wylogowuje się",
			FlagSet:     flag.NewFlagSet("xades-debug", flag.ExitOnError),
			Run:         xadesDebugRun,
		},
	}

	XadesDebugCommand.FlagSet.StringVar(&xadesDebugArgs.signedFile, "f", "", "ścieżka do podpisanego pliku (np. AuthTokenRequest.signed.xml)")
	testGatewayFlag(XadesDebugCommand.FlagSet)

	registerCommand(&XadesDebugCommand.Command)
}

func xadesDebugRun(c *Command) error {
	var cfg = config.GetConfig()

	authValidator := esignature.NewESignatureTokenHandler(
		cfg.APIConfig(environment),
		"",
		xadesDebugArgs.signedFile,
	)
	cli, err := v2.NewClient(c.Context, cfg, environment, v2.WithAuthValidator(authValidator))
	if err != nil {
		return err
	}
	if err := cli.ObtainToken(); err != nil {
		return err
	}

	return cli.Logout()
}
