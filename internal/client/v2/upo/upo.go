package upo

import (
	"context"
	"fmt"
	HTTP "ksef/internal/http"
	"ksef/internal/utils"
	"os"
	"path/filepath"
	"time"
)

type UPODownloadPage struct {
	ReferenceNumber string `json:"referenceNumber" yaml:"ref-no"`
	DownloadUrl     string `json:"downloadUrl" yaml:"downloadUrl"`
}

type UPODownloaderParams struct {
	Path  string
	Mkdir bool
	Wait  time.Duration
}

// prawdopodobny adres:
// https://ksef.mf.gov.pl/web/anonymous-upo-status
// const downloadUrl = "https://klient-eformularz.mf.gov.pl/api/upo/%s/pdf"
type UPODownloader struct {
	httpClient *HTTP.Client
	params     UPODownloaderParams
}

func NewDownloader(httpClient *HTTP.Client, params UPODownloaderParams) *UPODownloader {
	return &UPODownloader{
		httpClient: httpClient,
		params:     params,
	}
}

func (ud *UPODownloader) Download(
	ctx context.Context,
	uploadSessionId string,
	pages []UPODownloadPage,
	success func(upoXMLFilename string),
) error {
	// outputPath
	_, err := utils.ResolveFilepath(
		utils.FilepathResolverConfig{
			Path:  ud.params.Path,
			Mkdir: ud.params.Mkdir,
			DefaultFilename: fmt.Sprintf(
				"%s.xml",
				uploadSessionId,
			),
		},
	)

	if err == utils.ErrDoesNotExistAndMkdirNotSpecified {
		return fmt.Errorf("wskazany katalog nie istnieje a nie użyłeś opcji `-m`")
	}
	if err != nil {
		return fmt.Errorf("błąd tworzenia katalogu wyjściowego: %v", err)
	}

	for _, page := range pages {
		outputFileName := filepath.Join(ud.params.Path, page.ReferenceNumber) + ".xml"
		outputFile, err := os.Create(outputFileName)
		if err != nil {
			return err
		}
		if err := ud.httpClient.Download(
			ctx,
			page.DownloadUrl,
			outputFile,
		); err != nil {
			return err
		}
		outputFile.Close()

		success(outputFileName)
	}

	return nil
}
