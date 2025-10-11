package xades

import (
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
	signCommand.Flags().StringVarP(&certFile, "cert", "", "", "plik certyfikatu")
	signCommand.Flags().StringVarP(&outputFile, "dest", "o", "AuthTokenRequest.signed.xml", "plik docelowy")

	signCommand.MarkFlagRequired("challenge")
	signCommand.MarkFlagRequired("cert")
	signCommand.MarkFlagRequired("dest")

	XadesCommand.AddCommand(signCommand)
}

func signChallengeFile(cmd *cobra.Command, _ []string) error {
	destFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer destFile.Close()
	return xades.SignAuthChallenge(challengeFile, certFile, destFile)
}
