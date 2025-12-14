package uploader

import (
	v2 "ksef/internal/client/v2"
	sessionTypes "ksef/internal/client/v2/session/types"
	"ksef/internal/invoicesdb/uploader/config"

	"github.com/spf13/viper"
)

type Uploader struct {
	Queue sessionTypes.UploadPayload
	// internal pointers
	ksefClient *v2.APIClient
	config     config.UploaderConfig
}

func NewUploader(vip *viper.Viper, config config.UploaderConfig, ksefClient *v2.APIClient) *Uploader {
	return &Uploader{
		Queue:      make(sessionTypes.UploadPayload),
		config:     config,
		ksefClient: ksefClient,
	}
}
