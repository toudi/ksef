package flags

import "github.com/spf13/cobra"

const (
	FlagNamePESEL = "pesel"
)

func PESEL(cmd *cobra.Command) {
	cmd.Flags().StringP(FlagNamePESEL, "p", "", "numer PESEL")
}
