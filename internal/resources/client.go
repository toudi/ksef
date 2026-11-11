package resources

import (
	"context"

	"ksef/core/client"
	"ksef/core/ksef"
	"ksef/internal/client/v2/auth/token"
	"ksef/internal/http"
	"ksef/internal/runtime"

	"github.com/spf13/viper"
)

var clientInstance map[string]*client.Client

func init() {
	clientInstance = make(map[string]*client.Client)
}

func HTTPClient(vip *viper.Viper) (*http.Client, error) {
	gateway, err := runtime.GatewayConfig(vip)
	if err != nil {
		return nil, err
	}

	return http.NewClient(gateway.APIBaseURL), nil
}

func Client(ctx context.Context, vip *viper.Viper, NIP string) (*client.Client, error) {
	mu.Lock()
	defer mu.Unlock()

	return clientLocked(ctx, vip, NIP)
}

func clientLocked(ctx context.Context, vip *viper.Viper, NIP string) (*client.Client, error) {
	if instance, exists := clientInstance[NIP]; exists {
		return instance, nil
	}

	keyring, err := Keyring(vip)
	if err != nil {
		return nil, err
	}

	certsDB, err := CertsDB(vip)
	if err != nil {
		return nil, err
	}

	invoiceDB, err := invoiceDBLocked(vip, NIP)
	if err != nil {
		return nil, err
	}

	dataStorage := invoiceDataStorageLocked(vip)

	gateway, err := runtime.GatewayConfig(vip)
	if err != nil {
		return nil, err
	}

	rateLimitsState, err := RateLimitsState(vip, NIP)
	if err != nil {
		return nil, err
	}

	authValidator := token.NewAuthHandler(
		token.WithNIP(NIP),
		token.WithCertsDB(certsDB),
		token.WithKeyring(keyring),
	)

	cli, err := client.Init(
		ctx,
		certsDB,
		invoiceDB,
		keyring,
		gateway,
		NIP,
		client.HTTPConfig{
			HTTPConfig: ksef.HTTPConfig{
				MaxRetries:   runtime.HttpMaxRetries(vip),
				PollWaitTime: runtime.HttpPollWaitTime(vip),
			},
			RateLimitsState: rateLimitsState,
		},
		client.WithAuthValidator(authValidator),
		client.WithPDFRenderer(NewPDFRenderer(vip, dataStorage)),
	)
	if err != nil {
		return nil, err
	}

	OnClose(func(c context.Context, _ error) error {
		return cli.Close()
	})

	clientInstance[NIP] = cli
	return cli, nil
}
