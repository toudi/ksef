package authorization

import (
	"ksef/cmd/ksef/commands/client"
	"ksef/cmd/ksef/flags"
	"ksef/internal/client/v2/auth"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	vip := viper.GetViper()
	vip.Set(auth.FlagExitAfterPersistingToken, "true")

	cli, err := client.InitClient(cmd)
	if err != nil {
		return err
	}
	if err := cli.ObtainToken(); err != nil {
		return err
	}
	if err = cli.WaitForTokenManagerLoop(); err != nil {
		return err
	}
	return cli.Logout()
}
