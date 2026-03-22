package invoicesdb

import (
	"ksef/internal/invoicesdb/config"
	"ksef/internal/runtime"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func GetAllNIPs(vip *viper.Viper) ([]string, error) {
	cfg := config.GetInvoicesDBConfig(vip)
	environmentId := runtime.GetEnvironmentId(vip)

	var nips []string
	root := filepath.Join(
		cfg.Root,
		environmentId,
	)

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			nips = append(nips, entry.Name())
		}
	}

	return nips, err
}
