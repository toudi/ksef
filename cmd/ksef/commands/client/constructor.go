package client

import (
	"ksef/core/certsdb"
	"ksef/core/certsdb/fs"
	kr "ksef/core/keyring"
	v2 "ksef/internal/client/v2"
	"ksef/internal/client/v2/auth/token"
	"ksef/internal/runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitClient(
	cmd *cobra.Command,
	vip *viper.Viper,
	keyring kr.Keyring,
	initializers ...v2.InitializerFunc,
) (*v2.APIClient, error) {
	var err error

	var cli *v2.APIClient

	cmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		// TODO: handle logout parameter here
		return nil
	}

	certsDB, err := CertsDB(vip)
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

func CertsDB(vip *viper.Viper) (*certsdb.CertsDB, error) {
	certsStorage, err := fs.Init(vip)
	if err != nil {
		return nil, err
	}
	return certsdb.Init(
		certsStorage,
		certsdb.WithEnvironment(runtime.GetEnvironmentId(vip)),
		certsdb.WithPreferredCertProfile(runtime.GetPreferredCertProfile(vip)),
	)
}
