package flags

import (
	"github.com/spf13/pflag"
)

const (
	FlagNameNIP = "nip"
)

func NIP(flagSet *pflag.FlagSet) {
	flagSet.StringP(FlagNameNIP, "n", "", "numer NIP podmiotu")
}
