package subjectsettings

import (
	"ksef/cmd/ksef/flags"

	"github.com/spf13/cobra"
)

var SubjectSettings = &cobra.Command{
	Use:   "subject-settings",
	Short: "zarządzanie ustawieniami podmiotu",
}

func init() {
	flags.NIP(SubjectSettings.PersistentFlags())
	SubjectSettings.MarkFlagRequired(flags.FlagNameNIP)
	SubjectSettings.AddCommand(copyPDFRendererConfig)
}
