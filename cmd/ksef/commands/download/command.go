package download

import (
	"fmt"

	"github.com/spf13/cobra"
)

var DownloadCommand = &cobra.Command{
	Use:   "download",
	Short: "pobiera faktury z KSeF",
	RunE:  downloadRun,
}

func init() {
	flags := DownloadCommand.Flags()
	registerFlags(flags)
}

func downloadRun(cmd *cobra.Command, _ []string) error {
	params, err := getDownloadParams(cmd.Flags())
	if err != nil {
		return err
	}

	fmt.Printf("params: %+v\n", params)
	return nil
}
