package commands

import (
	"errors"
	"fmt"
	"ksef/cmd/ksef/commands/client"
	v2 "ksef/internal/client/v2"
	"ksef/internal/client/v2/upo"
	"ksef/internal/logging"
	registryPkg "ksef/internal/registry"
	"ksef/internal/runtime"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagNameRegistry = "registry"
)

var statusCommand = &cobra.Command{
	Use:   "status",
	Short: "sprawdza status wysłanych faktur i pobiera dokument UPO",
	RunE:  statusRun,
}

var registryPath string

var upoDownloaderParams upo.UPODownloaderParams

func init() {
	flagSet := statusCommand.Flags()
	cobra.MarkFlagRequired(flagSet, flagNameRegistry)

	flagSet.StringVarP(&registryPath, flagNameRegistry, "r", "", "ścieżka do katalogu z rejestrem")
	flagSet.StringVarP(&upoDownloaderParams.Path, "output", "o", "", "ścieżka do zapisu UPO (domyślnie katalog rejestru + {nrRef}.pdf)")
	flagSet.BoolVarP(&upoDownloaderParams.Mkdir, "m", "", false, "stwórz ktalog do zapisu, jeśli wskazany nie istnieje")
	// flagSet.BoolFunc("xml", "zapis UPO jako plik XML", func(s string) error {
	// 	upoDownloaderParams.Format = upo.UPOFormatXML
	// 	return nil
	// })
	flagSet.DurationP("wait", "w", time.Duration(0), "czekaj na przetworzenie sesji (tryb synchroniczny)")
	flagSet.Lookup("wait").NoOptDefVal = "5m"

	flagSet.SortFlags = false
}

func statusRun(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()
	registry, err := registryPkg.LoadRegistry(registryPath)
	if err != nil {
		return fmt.Errorf("unable to load status from file: %v", err)
	}

	if registry.Environment == "" || registry.UploadSessions == nil {
		return fmt.Errorf(
			"file deserialized correctly, but either environment or referenceNo are empty: %+v",
			registry,
		)
	}

	if parsedDur, _ := cmd.Flags().GetDuration("wait"); parsedDur > 0 {
		upoDownloaderParams.Wait = parsedDur
	}

	runtime.SetGateway(vip, registry.Environment)

	if registry.Issuer == "" {
		return errors.New("nie znaleziono numeru NIP")
	}

	if upoDownloaderParams.Path == "" {
		upoDownloaderParams.Path = registry.Dir
	}

	var statusErr error

	defer func() {
		if statusErr == nil {
			logging.UploadLogger.Info("sprawdzanie statusu zakończone sukcesem, zapisuję plik rejestru")
			registry.Save("")
		} else {
			logging.UploadLogger.Error("sprawdzanie statusu zakończone niepowodzeniem, nie zapisuję zmian w rejestrze", "err", statusErr)
		}
	}()

	runtime.SetNIP(vip, registry.Issuer)

	cli, err := client.InitClient(cmd, v2.WithRegistry(registry))
	if err != nil {
		return fmt.Errorf("błąd inicjalizacji klienta: %v", err)
	}

	return cli.UploadSessionsStatusCheck(cmd.Context(), upoDownloaderParams)
}
