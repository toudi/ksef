package uploader

import (
	"context"
	"errors"
	v2 "ksef/internal/client/v2"
	"ksef/internal/invoicesdb/uploader/config"
	"time"
)

var (
	errTimeoutWaitingForStatus = errors.New("timed out waiting for upload session status")
	errCheckingStatus          = errors.New("error checking upload session status")
)

type UploadStatus struct {
	KSeFRefNo string
	Errors    []string
}

type UploadResult struct {
	// first map key is is the session ID.
	// value of the first map is another map, where the key is the invoice checksum
	// and the value is the upload status.
	UploadSessions map[string]map[string]*UploadStatus
}

func (u *Uploader) WaitForResult(ctx context.Context, params config.UploaderConfig, client *v2.APIClient) (*UploadResult, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, params.WaitTimeout)
	defer cancel()

	pollingTicker := time.NewTicker(5 * time.Second)
	defer pollingTicker.Stop()

	var finished bool = false
	var err error

	var result = &UploadResult{
		UploadSessions: make(map[string]map[string]*UploadStatus),
	}

	for !finished {
		select {
		case <-timeoutCtx.Done():
			return nil, errTimeoutWaitingForStatus
		case <-pollingTicker.C:
			finished, err = u.checkUploadSessionStatus(ctx, result)
			if err != nil {
				return nil, errors.Join(errCheckingStatus, err)
			}
		}
	}

	return result, nil
}
