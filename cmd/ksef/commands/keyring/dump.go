package keyring

import (
	"ksef/cmd/ksef/flags"
	"ksef/core/certsdb"
	"ksef/core/keyring"
	"ksef/internal/config"
	"ksef/internal/resources"
	"ksef/internal/runtime"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dumpKeyringToFileCommand = &cobra.Command{
	Use:     "dump",
	Short:   "zrzuca zawartość keyringu do zaszyfrowanego pliku",
	RunE:    dumpKeyringToFileRun,
	PreRunE: setKeyringKeys,
}

var allKeyringKeys []string

func init() {
	config.FileKeyringFlags(dumpKeyringToFileCommand.Flags())
	dumpKeyringToFileCommand.MarkFlagRequired(flags.FlagNameNIP)
}

func setKeyringKeys(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()
	ctx := cmd.Context()

	certsDB, err := resources.CertsDB(vip)
	if err != nil {
		return err
	}
	defer resources.Close(ctx)
	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return err
	}

	// iterate over all potential authorization keys to build keyring keys
	filtered, err := certsDB.FilterCertificates(ctx, func(c certsdb.Certificate) bool {
		return slices.Contains(c.NIP, nip) && slices.Contains(c.Usage, certsdb.CertificateUsageAuthentication)
	})
	if err != nil {
		return err
	}
	for _, cert := range filtered {
		allKeyringKeys = append(allKeyringKeys, keyring.SessionTokensKey(cert.UID))
	}
	return nil
}

func dumpKeyringToFileRun(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()

	environmentId := runtime.GetEnvironmentId(vip)
	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return err
	}

	fileKeyringCfg, err := keyring.GetFileBasedKeyringConfig(vip)
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
		systemKeyringValue, err := systemKeyring.Get(environmentId, nip, key)
		if err != nil {
			return err
		}
		if err = fileKeyring.Set(environmentId, nip, key, systemKeyringValue); err != nil {
			return err
		}
	}

	return fileKeyring.Close()
}
