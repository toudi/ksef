package authorization

import (
	"bytes"
	"ksef/cmd/ksef/commands/authorization/challenge"
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/auth/xades"
	"os"

	kr "ksef/internal/keyring"

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
	certsDB, err := certsdb.OpenOrCreate(vip)
	if err != nil {
		return err
	}
	keyring, err := kr.NewKeyring(vip)
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
	return xades.SignAuthChallenge(bytes.NewBuffer(challengeBytes), certFile, keyring, destFile)
}
