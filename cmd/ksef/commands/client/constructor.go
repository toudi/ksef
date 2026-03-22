package client

import (
	"ksef/internal/certsdb"
	v2 "ksef/internal/client/v2"
	"ksef/internal/client/v2/auth/token"
	kr "ksef/internal/keyring"

	kr "ksef/internal/keyring"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitClient(cmd *cobra.Command, vip *viper.Viper, keyring kr.Keyring, initializers ...v2.InitializerFunc) (*v2.APIClient, error) {
	var err error

	var cli *v2.APIClient

	cmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		// TODO: handle logout parameter here
		return nil
	}
	certsDB, err := certsdb.OpenOrCreate(vip)
	if err != nil {
		return nil, err
	}

	keyring, err := kr.NewKeyring(vip)
	if err != nil {
		return nil, err
	}

	clientInitializers := []v2.InitializerFunc{
		v2.WithKeyring(keyring),
		v2.WithAuthValidator(
			token.NewAuthHandler(
				vip,
				token.WithCertsDB(certsDB),
				token.WithKeyring(keyring),
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
