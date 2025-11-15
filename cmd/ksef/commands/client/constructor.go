package client

import (
	"ksef/internal/certsdb"
	v2 "ksef/internal/client/v2"
	"ksef/internal/client/v2/auth/token"
	"ksef/internal/runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitClient(cmd *cobra.Command, initializers ...v2.InitializerFunc) (*v2.APIClient, error) {
	vip := viper.GetViper()
	var err error
	var env = runtime.GetGateway(vip)
	var cli *v2.APIClient

	cmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		// TODO: handle logout parameter here
		return nil
	}
	certsDB, err := certsdb.OpenOrCreate(env)
	if err != nil {
		return nil, err
	}

	var clientInitializers = []v2.InitializerFunc{
		v2.WithAuthValidator(
			token.NewAuthHandler(
				vip,
				token.WithCertsDB(certsDB),
			),
		),
		v2.WithCertificatesDB(certsDB),
	}

	clientInitializers = append(clientInitializers, initializers...)

	cli, err = v2.NewClient(
		cmd.Context(),
		vip,
		clientInitializers...,
	)

	return cli, err
}
