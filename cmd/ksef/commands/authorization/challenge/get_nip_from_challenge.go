package challenge

import (
	"encoding/xml"
	"ksef/internal/config"
	"os"

	"github.com/spf13/viper"
)

type AuthChallenge struct {
	XMLName xml.Name `xml:"AuthTokenRequest"`
	NIP     string   `xml:"ContextIdentifier>Nip"`
}

func GetNIPFromChallengeFile(challengeFile string) (challengeBytes []byte, nip string, err error) {
	challengeBytes, err = os.ReadFile(challengeFile)
	if err != nil {
		return nil, "", err
	}
	var challenge AuthChallenge
	if err = xml.Unmarshal(challengeBytes, &challenge); err != nil {
		return nil, "", err
	}
	config.SetNIP(viper.GetViper(), challenge.NIP)
	return challengeBytes, challenge.NIP, nil
}
