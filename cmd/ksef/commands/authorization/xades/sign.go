package xades

import (
	"bytes"
	"ksef/cmd/ksef/commands/authorization/challenge"
	"ksef/internal/certsdb"
	"ksef/internal/environment"
	"ksef/internal/sei/api/client/v2/auth/xades"
	"os"

	"github.com/spf13/cobra"
)

var signCommand = &cobra.Command{
	Use:   "sign",
	Short: "podpisuje wskazany plik wyzwania przy użyciu certyfikatu",
	RunE:  signChallengeFile,
}

var outputFile string

func init() {
	signCommand.Flags().StringVarP(&challengeFile, "challenge", "f", "", "plik wyzwania")
	signCommand.Flags().StringVarP(&outputFile, "dest", "o", "AuthTokenRequest.signed.xml", "plik docelowy")

	signCommand.MarkFlagRequired("challenge")

	XadesCommand.AddCommand(signCommand)
}

func signChallengeFile(cmd *cobra.Command, _ []string) error {
	env := environment.FromContext(cmd.Context())
	certsDB, err := certsdb.OpenOrCreate(env)
	if err != nil {
		return err
	}
	challengeBytes, nip, err := challenge.GetNIPFromChallengeFile(challengeFile)
	if err != nil {
		return err
	}
	certFile, err := certsDB.GetByUsage(certsdb.UsageAuthentication, nip)
	if err != nil {
		return err
	}

	destFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer destFile.Close()
	return xades.SignAuthChallenge(bytes.NewBuffer(challengeBytes), certFile, destFile)
}
