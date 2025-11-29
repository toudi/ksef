package download

import (
	"ksef/cmd/ksef/commands/client"
	v2 "ksef/internal/client/v2"
	"ksef/internal/client/v2/invoices"
	"ksef/internal/logging"
	registryPkg "ksef/internal/registry"

	"github.com/spf13/cobra"
)

var DownloadCommand = &cobra.Command{
	Use:   "download [registry-dir]",
	Short: "pobiera faktury z KSeF do wskazanego katalogu rejestru lub odświeża istniejący",
	RunE:  downloadRun,
	Args:  cobra.ExactArgs(1),
}

func init() {
	flags := DownloadCommand.Flags()
	registerFlags(flags)
}

func downloadRun(cmd *cobra.Command, args []string) error {
	params, err := getDownloadParams(cmd.Flags())
	if err != nil {
		return err
	}

	registryDir := args[0]
	registry, err := registryPkg.OpenOrCreate(registryDir)
	if err != nil {
		return err
	}
	cli, err := client.InitClient(
		cmd, v2.WithRegistry(registry),
	)
	for _, subjectType := range params.SubjectTypes {
		logging.DownloadLogger.Info("pobieram faktury dla typu", "subjectType", subjectType)
		if err = cli.SyncInvoices(
			cmd.Context(),
			invoices.SyncParams{
				DestPath:      registryDir,
				SubjectType:   subjectType,
				PageSize:      params.PageSize,
				DateRangeType: params.DateType,
				DateFrom:      params.StartDate,
			},
		); err != nil {
			return err
		}
	}
	return nil
}
