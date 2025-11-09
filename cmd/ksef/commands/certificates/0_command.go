package certificates

import (
	v2 "ksef/internal/client/v2"

	"github.com/spf13/cobra"
)

var CertificatesCommand = &cobra.Command{
	Use:   "certs",
	Short: "zarządzanie certyfikatami",
}

var cli *v2.APIClient
var err error
