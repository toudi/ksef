package client

import (
	kr "ksef/core/keyring"
	v2 "ksef/internal/client/v2"
	"ksef/internal/client/v2/auth/token"
	"ksef/internal/resources"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitClient(
	cmd *cobra.Command,
	vip *viper.Viper,
	keyring kr.Keyring,
	initializers ...v2.InitializerFunc,
) (*v2.APIClient, error) {
	certsDB, err := resources.CertsDB(vip)
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

	return v2.NewClient(
		cmd.Context(),
		vip,
		clientInitializers...,
	)
}
