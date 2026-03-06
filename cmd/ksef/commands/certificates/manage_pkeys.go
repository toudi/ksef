package certificates

import (
	"ksef/internal/certsdb"
	"ksef/internal/logging"

	kr "ksef/internal/keyring"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pkeysManagementCommand = &cobra.Command{
	Use:   "pkeys",
	Short: "zarządzanie kluczami prywatnymi",
	RunE:  cobra.NoArgs,
}

var pkeysEncryptCommand = &cobra.Command{
	Use:   "encrypt",
	Short: "zaszyfruj klucze prywatne zapisując hasło do keyringu",
	RunE:  pkeysEncrypt,
}

func init() {
	pkeysManagementCommand.AddCommand(
		pkeysEncryptCommand,
	)
	CertificatesCommand.AddCommand(
		pkeysManagementCommand,
	)
}

func pkeysEncrypt(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()

	certsDB, err := certsdb.OpenOrCreate(vip)
	if err != nil {
		return err
	}

	var keyring kr.Keyring // TODO: initialize / save

	if err := certsDB.ForEach(
		func(cert certsdb.Certificate) bool { return true },
		func(cert *certsdb.Certificate) error {
			isEncrypted, err := cert.PrivateKeyIsEncrypted()
			if err != nil {
				return err
			}
			if isEncrypted {
				logging.CertsDBLogger.Info("klucz prywatny certyfikatu jest już zaszyfrowany. no-op", "cert id", cert.UID)
				return nil
			}
			return cert.EncryptPrivateKey(keyring)
		},
	); err != nil {
		return err
	}

	return certsDB.Save()
}
