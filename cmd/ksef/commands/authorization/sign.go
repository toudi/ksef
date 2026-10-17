package authorization

import (
	"bytes"
	"ksef/cmd/ksef/commands/authorization/challenge"
	"ksef/cmd/ksef/commands/client"
	"ksef/core/certsdb"
	"ksef/internal/client/v2/auth/xades"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var signCommand = &cobra.Command{
	Use:   "sign <challenge-file>",
	Short: "podpisuje wskazany plik wyzwania przy użyciu certyfikatu",
	RunE:  signChallengeFile,
	Args:  cobra.ExactArgs(1),
}

var outputFile string

func init() {
	signCommand.Flags().StringVarP(&outputFile, "dest", "o", "AuthTokenRequest.signed.xml", "plik docelowy")

	AuthCommand.AddCommand(signCommand)
}

func signChallengeFile(cmd *cobra.Command, args []string) error {
	vip := viper.GetViper()
	challengeFile = args[0]
	certsDB, err := client.CertsDB(vip)
	if err != nil {
		return err
	}
	challengeBytes, nip, err := challenge.GetNIPFromChallengeFile(challengeFile)
	if err != nil {
		return err
	}
	cert, err := certsDB.GetCertificateByUsageAndNIP(
		cmd.Context(), certsdb.CertificateUsageAuthentication, nip,
	)
	if err != nil {
		return err
	}

	destFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer destFile.Close()
	return xades.SignAuthChallenge(
		cmd.Context(),
		bytes.NewBuffer(challengeBytes),
		certsDB,
		cert.UID,
		destFile,
	)
}
