package commands

import (
	"flag"
	"ksef/internal/config"
	v2 "ksef/internal/sei/api/client/v2"
	"ksef/internal/sei/api/client/v2/auth/esignature"
)

type xadesInitCommand struct {
	Command
}

type xadesArgsType struct {
	nip string
}

var XadesInitCommand *xadesInitCommand
var xadesArgs xadesArgsType

func init() {
	XadesInitCommand = &xadesInitCommand{
		Command: Command{
			Name:        "xades-init",
			Description: "generuje plik autoryzacji xml do podpisania (np. przez system ePUAP)",
			FlagSet:     flag.NewFlagSet("xades-init", flag.ExitOnError),
			Run:         xadesInitRun,
		},
	}

	XadesInitCommand.FlagSet.StringVar(&xadesArgs.nip, "nip", "", "numer NIP podmiotu")
	testGatewayFlag(XadesInitCommand.FlagSet)

	registerCommand(&XadesInitCommand.Command)
}

func xadesInitRun(c *Command) error {
	var cfg = config.GetConfig()

	authValidator := esignature.NewESignatureTokenHandler(
		cfg.APIConfig(environment),
		xadesArgs.nip,
		"",
	)
	cli, err := v2.NewClient(c.Context, cfg, environment, v2.WithAuthValidator(authValidator))
	if err != nil {
		return err
	}
	return cli.WaitForTokenManagerLoop()
}
