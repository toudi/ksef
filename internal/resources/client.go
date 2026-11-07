package resources

import (
	"ksef/core/client"
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

func Client(vip *viper.Viper, NIP string) (*client.Client, error) {
	mu.Lock()
	defer mu.Unlock()

	return clientLocked(vip, NIP)
}

func clientLocked(vip *viper.Viper, NIP string) (*client.Client, error) {
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

	cli := client.Init(
		certsDB,
		invoiceDB,
		keyring,
		gateway,
		client.WithPDFRenderer(NewPDFRenderer(vip, dataStorage)),
	)

	clientInstance[NIP] = cli
	return cli, nil
}
