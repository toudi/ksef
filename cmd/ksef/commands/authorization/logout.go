package authorization

import (
	"ksef/cmd/ksef/commands/client"
	"ksef/cmd/ksef/flags"

	"github.com/spf13/cobra"
)

var logoutCommand = &cobra.Command{
	Use:   "logout",
	Short: "zamknij zachowaną sesję logowania",
	RunE:  logout,
}

func init() {
	flags.NIP(logoutCommand.Flags())
	_ = cobra.MarkFlagRequired(logoutCommand.Flags(), flags.FlagNameNIP)
	AuthCommand.AddCommand(logoutCommand)
}

func logout(cmd *cobra.Command, _ []string) error {
	cli, err := client.InitClient(cmd)
	if err != nil {
		return err
	}
	if err = cli.WaitForTokenManagerLoop(); err != nil {
		return err
	}
	return cli.Logout()
}
