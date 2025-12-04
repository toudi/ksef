package keyring

import (
	"ksef/cmd/ksef/flags"
	"ksef/internal/runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var KeyringCommand = &cobra.Command{
	Use:     "keyring",
	Short:   "zarządzanie keyringiem",
	PreRunE: setNip,
}

func init() {
	flags.NIP(KeyringCommand.PersistentFlags())
	KeyringCommand.AddCommand(dumpKeyringToFileCommand)
	KeyringCommand.AddCommand(loadKeyringFromFileCommand)
}

func setNip(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()
	return runtime.SetNIPFromFlags(vip, cmd.Flags())
}
