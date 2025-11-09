package client

import (
	"ksef/internal/certsdb"
	v2 "ksef/internal/client/v2"
	"ksef/internal/client/v2/auth/token"
	"ksef/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitClient(cmd *cobra.Command) (*v2.APIClient, error) {
	vip := viper.GetViper()
	var err error
	var env = config.GetGateway(vip)
	var cli *v2.APIClient

	nip, err := config.GetNIP(vip)
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
		config.GetGateway(vip),
		v2.WithAuthValidator(
			token.NewAuthHandler(config.GetGateway(vip), nip),
		),
		v2.WithCertificatesDB(certsDB),
	)

	return cli, err
}
