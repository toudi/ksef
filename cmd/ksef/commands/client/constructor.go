package client

import (
	"ksef/cmd/ksef/flags"
	"ksef/internal/certsdb"
	"ksef/internal/config"
	"ksef/internal/environment"
	v2 "ksef/internal/sei/api/client/v2"
	"ksef/internal/sei/api/client/v2/auth/token"

	"github.com/spf13/cobra"
)

func InitClient(cmd *cobra.Command) (*v2.APIClient, error) {
	var err error
	var env = environment.FromContext(cmd.Context())
	var cli *v2.APIClient

	nip, err := cmd.Flags().GetString(flags.FlagNameNIP)
	if err != nil {
		return nil, err
	}
	cmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		// TODO: handle logout parameter here
		return nil
	}
	certsDB, err := certsdb.OpenOrCreate(env)
	if err != nil {
		return nil, err
	}

	cli, err = v2.NewClient(
		cmd.Context(),
		config.GetConfig(),
		env,
		v2.WithAuthValidator(
			token.NewAuthHandler(config.GetConfig().APIConfig(env), nip),
		),
		v2.WithCertificatesDB(certsDB),
	)

	return cli, err
}
