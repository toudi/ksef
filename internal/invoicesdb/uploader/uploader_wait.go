package uploader

import (
	"context"
	"errors"
	"ksef/internal/client/v2/session/status"
	sessionTypes "ksef/internal/client/v2/session/types"
	"time"
)

var (
	errTimeoutWaitingForStatus = errors.New("timed out waiting for upload session status")
	errCheckingStatus          = errors.New("error checking upload session status")
)

func (u *Uploader) WaitForResult(
	ctx context.Context,
	result []*sessionTypes.UploadSessionResult,
) ([]*sessionTypes.UploadSessionResult, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, u.config.WaitTimeout)
	defer cancel()

	pollingTicker := time.NewTicker(5 * time.Second)
	defer pollingTicker.Stop()

	var finished bool = false
	var err error

	var statusChecker *status.SessionStatusChecker = u.ksefClient.SessionStatusChecker()

	for !finished {
		select {
		case <-timeoutCtx.Done():
			return nil, errTimeoutWaitingForStatus
		case <-pollingTicker.C:
			finished, err = u.checkUploadSessionStatus(ctx, result, statusChecker)
			if err != nil {
				return nil, errors.Join(errCheckingStatus, err)
			}
		}
	}

	return result, nil
}
