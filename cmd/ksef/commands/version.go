package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version string
	date    string
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Zwraca informację o wersji",
	RunE:  printVersionInfo,
}

func init() {
	RootCommand.AddCommand(versionCommand)
}

func printVersionInfo(cmd *cobra.Command, _ []string) error {
	fmt.Printf("klient KSeF %s\ndata wydania %s\n", version, date)
	return nil
}
