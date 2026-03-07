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

var pkeysDecryptCommand = &cobra.Command{
	Use:   "decrypt",
	Short: "odszyfruj klucze prywatne używając kluczy z keyringu",
	RunE:  pkeysDecrypt,
}

func init() {
	pkeysManagementCommand.AddCommand(
		pkeysEncryptCommand, pkeysDecryptCommand,
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

	keyring, err := kr.NewKeyring(vip)
	if err != nil {
		return err
	}
	// this is required, since if somebody would be using file-based keyring we
	// actually have to save it
	defer keyring.Close()

	if err := certsDB.ForEach(
		func(cert certsdb.Certificate) bool { return true },
		func(cert *certsdb.Certificate) error {
			if !cert.IsEncryptable() {
				return nil
			}
			isEncrypted, err := cert.PrivateKeyIsEncrypted()
			if err != nil {
				return err
			}
			if isEncrypted {
				logging.CertsDBLogger.Info("klucz prywatny certyfikatu jest już zaszyfrowany. no-op", "cert", cert.UID)
				return nil
			}
			logging.CertsDBLogger.Info("Szyfruję klucz prywatny certyfikatu", "cert", cert.UID)
			return cert.EncryptPrivateKey(keyring)
		},
	); err != nil {
		return err
	}

	return certsDB.Save()
}

func pkeysDecrypt(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()

	certsDB, err := certsdb.OpenOrCreate(vip)
	if err != nil {
		return err
	}

	keyring, err := kr.NewKeyring(vip)
	if err != nil {
		return err
	}

	if err := certsDB.ForEach(
		func(cert certsdb.Certificate) bool { return true },
		func(cert *certsdb.Certificate) error {
			if !cert.IsEncryptable() {
				return nil
			}
			isEncrypted, err := cert.PrivateKeyIsEncrypted()
			if err != nil {
				return err
			}
			if !isEncrypted {
				logging.CertsDBLogger.Info("klucz prywatny certyfikatu jest już odszyfrowany. no-op", "cert", cert.UID)
			}
			if isEncrypted {
				logging.CertsDBLogger.Info("Odszyfrowuję klucz prywatny certyfikatu", "cert", cert.UID)
				return cert.DecryptPrivateKey(keyring)
			}
			return nil
		},
	); err != nil {
		return err
	}

	return certsDB.Save()
}
