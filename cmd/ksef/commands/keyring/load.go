package keyring

import (
	"ksef/internal/config"
	"ksef/internal/keyring"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loadKeyringFromFileCommand = &cobra.Command{
	Use:   "load",
	Short: "przenosi wartości z zaszyfrowanego pliku do systemowego keyringu",
	RunE:  loadKeyringFromFileRun,
}

func init() {
	config.FileKeyringFlags(loadKeyringFromFileCommand.Flags())
}

func loadKeyringFromFileRun(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()

	fileKeyringCfg, err := config.GetFileBasedKeyringConfig(vip)
	if err != nil {
		return err
	}
	// there are bunch of keys so there's hardly any point re-saving the file each time
	fileKeyringCfg.Buffered = true
	fileKeyring, err := keyring.NewFileBasedKeyring(fileKeyringCfg)
	if err != nil {
		return err
	}
	systemKeyring := keyring.NewSystemKeyring()

	return fileKeyring.ForEach(func(bucket, nip, key, contents string) error {
		return systemKeyring.Set(bucket, nip, key, contents)
	})

}
