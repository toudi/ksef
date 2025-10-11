package xades

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"ksef/internal/config"
	"ksef/internal/environment"

	"github.com/spf13/cobra"
)

var (
	pesel                     string
	errInvalidNIP             = errors.New("nieprawidłowy numer NIP")
	errTestGatewayNotSelected = errors.New("komenda działa tylko dla bramki testowej")
)

var genSelfSignedCertCommand = &cobra.Command{
	Use:   "generate-cert",
	Short: "generuje samopodpisany certyfikat dla podanego numeru NIP (jedynie przy użyciu bramki testowej)",
	RunE:  generateSelfSignedCert,
}

func init() {
	genSelfSignedCertCommand.Flags().StringVarP(&pesel, "pesel", "p", "", "numer PESEL (w przypadku użycia tej flagi zostanie wygenerowany certyfikat dla osoby fizycznej)")
	XadesCommand.AddCommand(genSelfSignedCertCommand)
}

func generateSelfSignedCert(cmd *cobra.Command, _ []string) error {
	cfg := config.GetConfig()
	env := environment.FromContext(cmd.Context())

	if env != environment.Test {
		return errTestGatewayNotSelected
	}

	var certBasename = "individual"
	var subject pkix.Name
	subject.Country = []string{"PL"}

	if pesel != "" {
		certBasename += "-" + pesel
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
		nipValidator := cfg.APIConfig(env).Environment.NIPValidator
		nip, err := cmd.Flags().GetString("nip")
		if err != nil {
			return err
		}
		if !nipValidator(nip) {
			return errInvalidNIP
		}
		certBasename = "company-" + nip
		subject.Organization = []string{"Gżegżółka sp z.o.o."}
		subject.CommonName = "Gżegżółka"
		subject.ExtraNames = append(subject.ExtraNames, []pkix.AttributeTypeAndValue{
			{
				Type:  asn1.ObjectIdentifier{2, 5, 4, 97},
				Value: "VATPL-" + nip,
			},
		}...)
	}

	return cfg.APIConfig(env).CertificatesDB.GenerateSelfSignedCert(
		subject, certBasename,
	)
}
