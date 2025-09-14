package commands

import (
	"flag"
	"ksef/internal/config"
	client_v2 "ksef/internal/sei/api/client/v2"
)

type downloadCertsCommand struct {
	Command
}

type downloadCertsArgsType struct {
	test bool
}

var downloadCertsArgs downloadCertsArgsType

func init() {
	DownloadCertsCommand := &downloadCertsCommand{
		Command: Command{
			Name:        "download-certs",
			FlagSet:     flag.NewFlagSet("download-certs", flag.ExitOnError),
			Description: "pobiera certyfikaty klucza publicznego",
			Run:         downloadCertsRun,
		},
	}

	testGatewayFlag(DownloadCertsCommand.FlagSet)

	registerCommand(&DownloadCertsCommand.Command)
}

func downloadCertsRun(c *Command) error {
	apiClient, err := client_v2.NewClient(c.Context, config.GetConfig(), environment)
	if err != nil {
		return err
	}

	return apiClient.DownloadCertificates(c.Context)
}
