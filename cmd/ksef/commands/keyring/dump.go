package keyring

import (
	"ksef/cmd/ksef/flags"
	"ksef/internal/config"
	"ksef/internal/keyring"
	"ksef/internal/runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dumpKeyringToFileCommand = &cobra.Command{
	Use:   "dump",
	Short: "zrzuca zawartość keyringu do zaszyfrowanego pliku",
	RunE:  dumpKeyringToFileRun,
}

func init() {
	config.FileKeyringFlags(dumpKeyringToFileCommand.Flags())
	dumpKeyringToFileCommand.MarkFlagRequired(flags.FlagNameNIP)
}

var allKeyringKeys = []string{keyring.KeySessionTokens}

func dumpKeyringToFileRun(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()

	gateway := runtime.GetGateway(vip)
	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return err
	}

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

	for _, key := range allKeyringKeys {
		systemKeyringValue, err := systemKeyring.Get(string(gateway), nip, key)
		if err != nil {
			return err
		}
		if err = fileKeyring.Set(string(gateway), nip, key, systemKeyringValue); err != nil {
			return err
		}
	}

	return fileKeyring.Close()
}
