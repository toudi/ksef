package authorization

import (
	"bytes"
	"ksef/cmd/ksef/commands/authorization/challenge"
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/auth/xades"
	"ksef/internal/config"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	AuthCommand.AddCommand(signCommand)
}

func signChallengeFile(cmd *cobra.Command, _ []string) error {
	certsDB, err := certsdb.OpenOrCreate(config.GetGateway(viper.GetViper()))
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
