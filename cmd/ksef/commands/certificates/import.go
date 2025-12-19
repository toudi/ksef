package certificates

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"ksef/cmd/ksef/flags"
	"ksef/internal/certsdb"
	"ksef/internal/runtime"
	"ksef/internal/utils"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var importCertificate = &cobra.Command{
	Use:   "import",
	Short: "importuje certyfikat wygenerowany przez MCU",
	RunE:  importCertificateRun,
}

var (
	errImportingCertificate = errors.New("błąd importu certyfikatu")
	errEncryptedPrivateKey  = errors.New("klucz prywatny jest zaszyfrowany - należy go odszyfrować za pomocą openssl")
	errParsingECPrivateKey  = errors.New("błąd parsowania klucza krzywych eliptycznych")
	errCopyingPKeyFile      = errors.New("błąd kopiowania pliku klucza prywatnego")
	errCopyingCertFile      = errors.New("błąd kopiowania pliku certyfikatu")
)

const (
	flagNamePrivateKey  = "private-key"
	flagNameCertificate = "certificate"
	flagNameSerial      = "serial"
	flagNameUsage       = "usage"
)

func init() {
	flagSet := importCertificate.Flags()
	flagSet.StringP(flagNamePrivateKey, "p", "", "plik klucza prywatnego")
	flagSet.String(flagNameCertificate, "", "plik cerfyfikatu")
	flagSet.String(flagNameSerial, "", "numer seryjny certyfikatu")
	flagSet.String(flagNameUsage, "", "przeznaczenie certyfikatu")
	flags.NIP(flagSet)

	flagSet.SortFlags = false

	cobra.MarkFlagRequired(flagSet, flagNamePrivateKey)
	cobra.MarkFlagRequired(flagSet, flagNameCertificate)
	cobra.MarkFlagRequired(flagSet, flagNameSerial)
	cobra.MarkFlagRequired(flagSet, flagNameUsage)
	cobra.MarkFlagRequired(flagSet, flags.FlagNameNIP)

	CertificatesCommand.AddCommand(importCertificate)
}

func importCertificateRun(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()
	var env = runtime.GetGateway(vip)

	certsDB, err := certsdb.OpenOrCreate(env)
	if err != nil {
		return err
	}

	// check if the key is of type ECDSA - we don't support any other keys
	privateKeyFile, _ := cmd.Flags().GetString(flagNamePrivateKey)
	privateKeyBytes, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return err
	}
	privateKeyDer, _ := pem.Decode(privateKeyBytes)
	if privateKeyDer.Type == "ENCRYPTED PRIVATE KEY" {
		return errEncryptedPrivateKey
	}
	privKey, err := x509.ParsePKCS8PrivateKey(privateKeyDer.Bytes)
	if err != nil {
		return errors.Join(errParsingECPrivateKey, err)
	}
	if _, ok := privKey.(*ecdsa.PrivateKey); !ok {
		return errParsingECPrivateKey
	}
	certificateFile, _ := cmd.Flags().GetString(flagNameCertificate)
	certificateData, err := os.ReadFile(certificateFile)
	if err != nil {
		return err
	}
	certificateDer, _ := pem.Decode(certificateData)
	certificate, err := x509.ParseCertificate(certificateDer.Bytes)
	if err != nil {
		return err
	}
	serialNumber, _ := cmd.Flags().GetString(flagNameSerial)
	usageString, _ := cmd.Flags().GetString(flagNameUsage)
	usage, err := certsdb.ValidateUsage(usageString)
	if err != nil {
		return err
	}
	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return err
	}

	// prepare certificate hash
	certificateHash := certsdb.CertificateHash{
		Environment: env,
		Usage:       []certsdb.Usage{usage},
		ValidFrom:   certificate.NotBefore,
		ValidTo:     certificate.NotAfter,
	}

	if err = certsDB.AddIfHashNotFound(certificateHash, func(newCert *certsdb.Certificate) (err error) {
		if err = utils.CopyFile(
			privateKeyFile, newCert.PrivateKeyFilename(),
		); err != nil {
			return errCopyingPKeyFile
		}

		if err = utils.CopyFile(
			certificateFile, newCert.Filename(),
		); err != nil {
			return errCopyingCertFile
		}

		newCert.NIP = nip
		newCert.SerialNumber = serialNumber
		newCert.Usage = []certsdb.Usage{usage}
		newCert.ValidFrom = certificate.NotBefore
		newCert.ValidTo = certificate.NotAfter

		return nil
	}); err != nil {
		return err
	}

	return certsDB.Save()
}
