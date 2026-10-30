package resources

import (
	"ksef/core/client"
	"ksef/internal/http"
	"ksef/internal/runtime"
	"sync"

	"github.com/spf13/viper"
)

var (
	clientInstances map[string]*client.Client
	mu              sync.Mutex
)

func init() {
	clientInstances = make(map[string]*client.Client)
}

func HTTPClient(vip *viper.Viper) (*http.Client, error) {
	gateway, err := runtime.GatewayConfig(vip)
	if err != nil {
		return nil, err
	}

	return http.NewClient(gateway.API), nil
}

func Client(NIP string, vip *viper.Viper) (*client.Client, error) {
	mu.Lock()
	defer mu.Unlock()

	if instance, exists := clientInstances[NIP]; exists {
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

	invoiceDB, err := InvoiceDB(vip)
	if err != nil {
		return nil, err
	}

	httpClient, err := HTTPClient(vip)
	if err != nil {
		return nil, err
	}

	client := client.Init(
		certsDB,
		invoiceDB,
		keyring,
		httpClient,
	)

	clientInstances[NIP] = client
	return client, nil
}
