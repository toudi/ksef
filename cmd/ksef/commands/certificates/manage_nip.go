package certificates

import (
	"ksef/cmd/ksef/commands/client"
	"ksef/cmd/ksef/flags"
	"ksef/internal/runtime"

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
	ctx := cmd.Context()

	certID := vip.GetString(flagCertID)
	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return err
	}

	certsDB, err := client.CertsDB(vip)
	if err != nil {
		return err
	}

	if err = certsDB.AddNIPToCertificate(ctx, certID, nip); err != nil {
		return err
	}
	return certsDB.Close(ctx)
}

func dropNIPFromCertRun(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()
	ctx := cmd.Context()

	certID := vip.GetString(flagCertID)
	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return err
	}

	certsDB, err := client.CertsDB(vip)
	if err != nil {
		return err
	}

	if err = certsDB.RemoveNIPFromCertificate(ctx, certID, nip); err != nil {
		return err
	}
	return certsDB.Close(ctx)
}
