package subjectsettings

import (
	"errors"
	"ksef/internal/config"
	invoicesdbconfig "ksef/internal/invoicesdb/config"
	subjectsettings "ksef/internal/invoicesdb/subject-settings"
	"ksef/internal/runtime"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var copyPDFRendererConfig = &cobra.Command{
	Use:   "copy-pdf-config",
	Short: "Kopiuje globalne ustawienia wydruku PDF do ustawień podmiotu",
	RunE:  copyPDFRendererConfigRun,
}

func init() {
	invoicesdbconfig.InvoicesDBFlags(copyPDFRendererConfig.Flags())
}

var errNoConfigToCopy = errors.New("brak konfiguracji do skopiowania")

func copyPDFRendererConfigRun(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()
	cfg := invoicesdbconfig.GetInvoicesDBConfig(vip)
	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return err
	}

	ss, err := subjectsettings.OpenOrCreate(
		filepath.Join(
			cfg.Root, runtime.GetEnvironmentId(vip), nip,
		),
	)
	if err != nil {
		return err
	}

	pdfConfig, err := config.GetPDFPrinterConfig(vip)
	if err != nil {
		return err
	}

	engines := pdfConfig.GetEngines()
	if len(engines) == 0 {
		return errNoConfigToCopy
	}

	if err := ss.Modify(func(state *subjectsettings.SubjectSettings) error {
		state.PDF = engines
		return nil
	}); err != nil {
		return err
	}

	return ss.Save()
}
