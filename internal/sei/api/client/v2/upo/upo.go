package upo

import (
	"context"
	"fmt"
	HTTP "ksef/internal/http"
	"ksef/internal/utils"
)

type UPODownloadFormat string

const (
	UPOFormatPDF UPODownloadFormat = "pdf"
	UPOFormatXML UPODownloadFormat = "xml"
)

type UPODownloadPage struct {
	ReferenceNumber string `json:"referenceNumber"`
	DownloadUrl     string `json:"downloadUrl"`
}

type UPODownloaderParams struct {
	Path   string
	Mkdir  bool
	Format UPODownloadFormat
}

// prawdopodobny adres:
// https://ksef.mf.gov.pl/web/anonymous-upo-status
// const downloadUrl = "https://klient-eformularz.mf.gov.pl/api/upo/%s/pdf"
type UPODownloader struct {
	httpClient HTTP.Client
	params     UPODownloaderParams
}

func NewDownloader(httpClient HTTP.Client, params UPODownloaderParams) *UPODownloader {
	return &UPODownloader{
		httpClient: httpClient,
		params:     params,
	}
}

func (ud *UPODownloader) Download(ctx context.Context, uploadSessionId string, pages []UPODownloadPage) error {
	outputPath, err := utils.ResolveFilepath(
		utils.FilepathResolverConfig{
			Path:  ud.params.Path,
			Mkdir: ud.params.Mkdir,
			DefaultFilename: fmt.Sprintf(
				"%s.%s",
				uploadSessionId,
				ud.params.Format,
			),
		},
	)

	if err == utils.ErrDoesNotExistAndMkdirNotSpecified {
		return fmt.Errorf("wskazany katalog nie istnieje a nie użyłeś opcji `-m`")
	}
	if err != nil {
		return fmt.Errorf("błąd tworzenia katalogu wyjściowego: %v", err)
	}

	// TODO: actual downloading through httpClient

	return nil
}
