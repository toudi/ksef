package certificates

import (
	"ksef/internal/logging"
	"ksef/internal/resources"

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
	ctx := cmd.Context()

	metaStorage, err := resources.CertsMetaStorage(vip)
	if err != nil {
		return err
	}
	dataStorage, err := resources.CertsDataStorage(vip)
	if err != nil {
		return err
	}
	defer resources.Close(ctx)

	certs, err := metaStorage.Certificates(ctx)
	if err != nil {
		return err
	}

	// let's iterate over certs from the storage and check whether they need encryption.
	for _, cert := range certs {
		if !cert.IsEncryptable() {
			continue
		}

		isEncrypted, err := dataStorage.IsPrivateKeyEncrypted(ctx, cert.UID)
		if err != nil {
			return err
		}

		if isEncrypted {
			// private key is already encrypted, can no-op
			logging.CertsDBLogger.Info("klucz prywatny certyfikatu jest już zaszyfrowany. no-op", "cert", cert.UID)
			continue
		}

		logging.CertsDBLogger.Info("Szyfruję klucz prywatny certyfikatu", "cert", cert.UID)
		if err := dataStorage.EncryptPrivateKey(ctx, cert.UID); err != nil {
			return err
		}
	}

	return nil
}

func pkeysDecrypt(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()
	ctx := cmd.Context()

	meta, err := resources.CertsMetaStorage(vip)
	if err != nil {
		return err
	}
	data, err := resources.CertsDataStorage(vip)
	if err != nil {
		return err
	}
	defer resources.Close(ctx)

	certs, err := meta.Certificates(ctx)
	if err != nil {
		return err
	}

	// let's iterate over certs from the storage and check whether they need decryption.
	for _, cert := range certs {
		if !cert.IsEncryptable() {
			continue
		}

		isEncrypted, err := data.IsPrivateKeyEncrypted(ctx, cert.UID)
		if err != nil {
			return err
		}

		if !isEncrypted {
			// private key is already decrypted, can no-op
			logging.CertsDBLogger.Info("klucz prywatny certyfikatu jest już odszyfrowany. no-op", "cert", cert.UID)
			continue
		}

		logging.CertsDBLogger.Info("Odszyfruję klucz prywatny certyfikatu", "cert", cert.UID)
		if err := data.DecryptPrivateKey(ctx, cert.UID); err != nil {
			return err
		}
	}

	return nil
}
