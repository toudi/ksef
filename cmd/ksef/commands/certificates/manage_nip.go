package certificates

import (
	"ksef/cmd/ksef/flags"
	"ksef/internal/certsdb"
	"ksef/internal/runtime"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addNIPToCert = &cobra.Command{
	Short: "Dodaje numer NIP do istniejącego certyfikatu",
	Use:   "add-nip",
	RunE:  addNIPToCertRun,
}

var dropNIPFromCert = &cobra.Command{
	Short: "Usuwa numer NIP z listy przypisanych do certyfikatu",
	Use:   "drop-nip",
	RunE:  dropNIPFromCertRun,
}

const (
	flagCertID = "cert-id"
)

func init() {
	flagSetAdd := addNIPToCert.Flags()
	flagSetDrop := dropNIPFromCert.Flags()
	flags.NIP(flagSetAdd)
	flags.NIP(flagSetDrop)
	flagSetAdd.String(flagCertID, "", "identyfikator certyfikatu")
	flagSetDrop.String(flagCertID, "", "identyfikator certyfikatu")
	addNIPToCert.MarkFlagRequired(flagCertID)
	addNIPToCert.MarkFlagRequired(flags.FlagNameNIP)
	dropNIPFromCert.MarkFlagRequired(flagCertID)
	dropNIPFromCert.MarkFlagRequired(flags.FlagNameNIP)

	CertificatesCommand.AddCommand(addNIPToCert, dropNIPFromCert)
}

func addNIPToCertRun(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()

	certID := vip.GetString(flagCertID)
	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return err
	}

	certDB, err := certsdb.OpenOrCreate(vip)
	if err != nil {
		return err
	}
	defer certDB.Save()

	return certDB.Upsert(
		func(cert certsdb.Certificate) bool {
			return cert.UID == certID
		},
		func(newCert *certsdb.Certificate) error {
			if !slices.Contains(newCert.NIP, nip) {
				newCert.NIP = append(newCert.NIP, nip)
			}
			return nil
		},
	)
}

func dropNIPFromCertRun(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()

	certID := vip.GetString(flagCertID)
	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return err
	}

	certDB, err := certsdb.OpenOrCreate(vip)
	if err != nil {
		return err
	}
	defer certDB.Save()

	return certDB.Upsert(
		func(cert certsdb.Certificate) bool {
			return cert.UID == certID
		},
		func(newCert *certsdb.Certificate) error {
			if slices.Contains(newCert.NIP, nip) {
				newCert.NIP = slices.DeleteFunc(newCert.NIP, func(n string) bool {
					return n == nip
				})
			}
			return nil
		},
	)
}
