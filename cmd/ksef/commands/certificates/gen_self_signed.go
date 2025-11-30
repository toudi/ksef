package certificates

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"ksef/cmd/ksef/flags"
	"ksef/internal/certsdb"
	"ksef/internal/runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var generateSelfSignedCommand = &cobra.Command{
	Use:     "gen-self-signed",
	Short:   "generuje samopodpisany certyfikat dla uwierzytelniania (*TYLKO* środowisko testowe)",
	RunE:    generateSelfSignedCert,
	PreRunE: validateParams,
}

var (
	errTestGatewayNotSelected = errors.New("komenda działa tylko dla bramki testowej")
	errInvalidNIP             = errors.New("nieprawidłowy numer NIP")
)

func init() {
	flags.PESEL(generateSelfSignedCommand)
	flags.NIP(generateSelfSignedCommand.Flags())

	CertificatesCommand.AddCommand(generateSelfSignedCommand)
}

func validateParams(cmd *cobra.Command, _ []string) error {
	env := runtime.GetGateway(viper.GetViper())

	if env != runtime.TestGateway {
		return errTestGatewayNotSelected
	}

	return nil
}

func generateSelfSignedCert(cmd *cobra.Command, _ []string) error {
	var env = runtime.GetGateway(viper.GetViper())
	certsDB, err := certsdb.OpenOrCreate(env)
	if err != nil {
		return err
	}
	defer certsDB.Save()

	pesel, err := cmd.Flags().GetString(flags.FlagNamePESEL)
	if err != nil {
		return err
	}

	var subject pkix.Name
	subject.Country = []string{"PL"}

	if pesel != "" {
		// we're generating a certificate for individual person
		subject.CommonName = "Michał Drzymała"
		subject.SerialNumber = "PESEL-" + pesel
		subject.ExtraNames = append(subject.ExtraNames, []pkix.AttributeTypeAndValue{
			{
				Type:  asn1.ObjectIdentifier{2, 5, 4, 42},
				Value: "Michał",
			},
			{
				Type:  asn1.ObjectIdentifier{2, 5, 4, 4},
				Value: "Drzymała",
			},
		}...)
	} else {
		// we're generating a certificate for a company, therefore NIP has to be set
		nip, err := runtime.GetNIP(viper.GetViper())
		if err != nil {
			return err
		}
		// certBasename = "company-" + nip
		subject.Organization = []string{"Gżegżółka sp z.o.o."}
		subject.CommonName = "Gżegżółka"
		subject.ExtraNames = append(subject.ExtraNames, []pkix.AttributeTypeAndValue{
			{
				Type:  asn1.ObjectIdentifier{2, 5, 4, 97},
				Value: "VATPL-" + nip,
			},
		}...)

	}

	return certsDB.GenerateSelfSignedCert(subject)
}
