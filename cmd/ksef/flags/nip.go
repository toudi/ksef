package flags

import (
	"ksef/internal/environment"

	"github.com/spf13/pflag"
)

const (
	FlagNameNIP = "nip"
)

func NIP(flagSet *pflag.FlagSet) {
	flagSet.StringP(FlagNameNIP, "n", "", "numer NIP podmiotu")
}

func GetNIP(flagSet *pflag.FlagSet, env environment.Environment) (string, error) {
	nip, err := flagSet.GetString(FlagNameNIP)
	if err != nil {
		return "", err
	}
	envConfig := environment.GetConfig(env)
	return nip, envConfig.NIPValidator(nip)
}
