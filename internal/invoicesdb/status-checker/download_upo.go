package statuschecker

import (
	"context"
	"ksef/internal/client/v2/session/types"
	"ksef/internal/client/v2/upo"
	"ksef/internal/http"
	"ksef/internal/invoicesdb/config"
	"ksef/internal/logging"
	"ksef/internal/pdf"
	"ksef/internal/pdf/printer"
	"ksef/internal/runtime"
	"path"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func getUPODownloadPath(vip *viper.Viper, month time.Time) (string, error) {
	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return "", err
	}

	environmentId := runtime.GetEnvironmentId(vip)
	invoicesDBConfig := config.GetInvoicesDBConfig(vip)

	return path.Join(
		invoicesDBConfig.Root,
		environmentId,
		nip,
		month.Format("2006"),
		month.Format("01"),
		"upo",
	), nil
}

func (c *StatusChecker) downloadUPO(
	ctx context.Context,
	uploadSession *types.UploadSessionResult,
) error {
	// during this pass of the Update function, we can finally assign the KSeF reference
	// numbers
	// dispatch upo downloader
	upoDestPath, err := getUPODownloadPath(c.vip, c.monthsRange[1])
	// c.monthsRange[1] is essentially today (i.e. this month)
	if err != nil {
		return err
	}

	var printer printer.PDFPrinter

	if c.cfg.UPODownloaderConfig.ConvertToPDF {
		printer, err = pdf.GetUPOPrinter(c.vip)
		if err != nil {
			return err
		}
	}

	upoDownloader := upo.NewDownloader(
		http.NewClient(""), upo.UPODownloaderParams{
			Path:  upoDestPath,
			Mkdir: true,
			Wait:  c.cfg.UPODownloaderConfig.Timeout,
		},
	)

	if err = upoDownloader.Download(
		ctx,
		uploadSession.SessionID,
		uploadSession.Status.Upo.Pages,
		func(upoXMLFilename string) {
			if c.cfg.UPODownloaderConfig.ConvertToPDF {
				if err = printer.PrintUPO(
					upoXMLFilename, strings.Replace(upoXMLFilename, ".xml", ".pdf", 1),
				); err != nil {
					logging.PDFRendererLogger.Error("błąd konwersji UPO do PDF", "src", upoXMLFilename, "err", err)
				}
			}
		},
	); err != nil {
		return err
	}

	return nil
}
