package security

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"ksef/internal/config"
	"ksef/internal/http"
	"os"
	"slices"
	"strings"
)

const endpointDownloadCertificates = "/api/v2/security/public-key-certificates"
const (
	usageKsefTokenEncryption    = "KsefTokenEncryption"
	usageSymmetricKeyEncryption = "SymmetricKeyEncryption"
)

type certificateRow struct {
	Certificate string   `json:"certificate"`
	Usage       []string `json:"usage"`
}

func DownloadCertificates(ctx context.Context, client *http.Client, cfg config.APIConfig) error {
	var certificates []certificateRow

	_, err := client.Request(
		ctx,
		http.RequestConfig{
			Dest:            &certificates,
			DestContentType: http.JSON,
		},
		endpointDownloadCertificates,
	)

	if err != nil {
		return err
	}

	for _, certificate := range certificates {
		if slices.Contains(certificate.Usage, usageKsefTokenEncryption) && slices.Contains(certificate.Usage, usageSymmetricKeyEncryption) {
			// so the base64-encoded content is actually the essense of PEM so let's use a nifty hack to save it
			pemFile, err := os.Create(cfg.Certificate.PEM())
			if err != nil {
				return err
			}
			defer pemFile.Close()
			if _, err = fmt.Fprintf(pemFile, "-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----\n", certificate.Certificate); err != nil {
				return err
			}

			// and the der file is just the binary version of it
			var base64Decoder = base64.NewDecoder(base64.StdEncoding, strings.NewReader(certificate.Certificate))
			derFile, err := os.Create(cfg.Certificate.DER())
			if err != nil {
				return err
			}
			defer derFile.Close()

			_, err = io.Copy(derFile, base64Decoder)
			if err != nil {
				return err
			}

			break
		}
	}

	return nil
}
