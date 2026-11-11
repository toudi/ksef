package client

import (
	"ksef/core/keyring"
	v2 "ksef/internal/client/v2"
	"ksef/internal/client/v2/auth/token"
	"ksef/internal/resources"
	"ksef/internal/runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitClient(
	cmd *cobra.Command,
	vip *viper.Viper,
	keyring keyring.Keyring,
	initializers ...v2.InitializerFunc,
) (*v2.APIClient, error) {
	certsDB, err := resources.CertsDB(vip)
	if err != nil {
		return nil, err
	}

	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return nil, err
	}

	rateLimitsState, err := resources.RateLimitsState(vip, nip)
	if err != nil {
		return nil, err
	}

	clientInitializers := []v2.InitializerFunc{
		v2.WithKeyring(keyring),
		v2.WithAuthValidator(
			token.NewAuthHandler(
				token.WithNIP(nip),
				token.WithCertsDB(certsDB),
				token.WithKeyring(keyring),
			),
		),
		v2.WithCertificatesDB(certsDB),
		v2.WithRateLimitsStateStorage(rateLimitsState),
	}

	clientInitializers = append(clientInitializers, initializers...)

	return v2.NewClient(
		cmd.Context(),
		vip,
		clientInitializers...,
	)
}
