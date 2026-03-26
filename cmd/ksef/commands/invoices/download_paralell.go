package invoices

import (
	"errors"
	"ksef/cmd/ksef/commands/client"
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/invoicesdb"
	downloaderconfig "ksef/internal/invoicesdb/downloader/config"
	kr "ksef/internal/keyring"
	"ksef/internal/logging"
	"ksef/internal/runtime"
	"slices"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var errOpeningCertsDB = errors.New("błąd otwierania bazy certyfikatów")

type downloadError struct {
	NIP string
	err error
}

func downloadRunParalell(cmd *cobra.Command, baseViper *viper.Viper, numWorkers int) error {
	// let's collect NIP's.
	// We can retrieve them from the certsDB and also pass the preferredCertProfile.
	// So effectively what we're doing is we're filtering available certs and for each
	// discovered cert we're iterating over the NIP numbers it handles.
	certsDB, err := certsdb.OpenOrCreate(baseViper)
	if err != nil {
		return errors.Join(errOpeningCertsDB, err)
	}

	var nipNumbers []string
	preferredProfile := runtime.GetPreferredCertProfile(baseViper)

	for cert := range certsDB.Filter(func(cert certsdb.Certificate) bool {
		if !slices.Contains(cert.Usage, certsdb.UsageAuthentication) {
			return false
		}
		return preferredProfile == "" || cert.MatchesProfile(preferredProfile)
	}) {
		for _, nip := range cert.NIP {
			if !slices.Contains(nipNumbers, nip) {
				nipNumbers = append(nipNumbers, nip)
			}
		}
	}

	downloaderConfig, err := downloaderconfig.GetDownloaderConfig(baseViper, "")
	if err != nil {
		return err
	}

	keyring, err := kr.NewKeyring(baseViper)
	if err != nil {
		logging.SeiLogger.Error("błąd inicjalizacji keyringu", "err", err)
		return err
	}
	defer keyring.Close()

	// now let's determine if there is less nip's to process than number of
	// declared workers. if so - let's decrement it to not waste resources
	if len(nipNumbers) < numWorkers {
		numWorkers = len(nipNumbers)
	}

	// start the workers
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	nipChannel := make(chan string, len(nipNumbers))
	errChannel := make(chan downloadError)

	for range numWorkers {
		go downloadWorker(cmd, &wg, baseViper, downloaderConfig, keyring, nipChannel, errChannel)
	}

	// start the error listening function
	var errors []downloadError
	go func() {
		for err := range errChannel {
			errors = append(errors, err)
		}
	}()

	// pass nip numbers to the workers
	for _, nip := range nipNumbers {
		nipChannel <- nip
	}

	// now we can close the channel. the worker will detect this as a signal that
	// there are no more incoming NIP's and that it can return
	close(nipChannel)

	// now we can wait for the workers to finish their work
	logging.DownloadLogger.Info("Oczekiwanie na zakończenie pobierania faktur")
	wg.Wait()
	close(errChannel)

	if len(errors) > 0 {
		logging.DownloadLogger.Info("Podczas pobierania wystąpiły następujące błędy")
		for _, err := range errors {
			logging.DownloadLogger.Error("Błąd pobierania", "NIP", err.NIP, "error", err)
		}
	}

	return nil
}

func downloadWorker(
	cmd *cobra.Command,
	wg *sync.WaitGroup,
	baseViper *viper.Viper,
	downloaderConfig invoices.DownloadParams,
	keyring kr.Keyring,
	nipChannel <-chan string,
	errorsChannel chan<- downloadError,
) {
	defer wg.Done()

	for nip := range nipChannel {
		vip := cloneViper(baseViper)
		runtime.SetNIP(vip, nip)

		if err := doDownload(cmd, vip, nip, downloaderConfig, keyring); err != nil {
			errorsChannel <- downloadError{
				NIP: nip,
				err: err,
			}
		}
	}
}

// this may seem weird, but it's actually easier to make this function look "almost" natural
// and return the error. This way, the "inner" function can focus on the actual logic of downloading invoices
// whereas the wrapper (i.e. worker in the worker pool) can focus on dealing with error channel and so on.
// from the perspective of this function - there's nothing magical - it just instantiates the KSeFClient
// and uses defer to close it.
func doDownload(
	cmd *cobra.Command,
	vip *viper.Viper,
	nip string,
	downloaderConfig invoices.DownloadParams,
	keyring kr.Keyring,
) error {
	ksefClient, err := client.InitClient(cmd, vip, keyring)
	if err != nil {
		return err
	}
	defer ksefClient.Close()

	invoicesDB, err := invoicesdb.OpenForNIP(
		nip, vip, invoicesdb.WithKSeFClient(ksefClient),
	)
	if err != nil {
		return err
	}

	return invoicesDB.DownloadInvoices(
		cmd.Context(), vip, downloaderConfig,
		logging.DownloadLogger.With("nip", nip),
	)
}
