package uploader

import (
	"context"
	"ksef/internal/client/v2/session/status"
	sessionTypes "ksef/internal/client/v2/session/types"
)

// checkUploadSessionStatus iterates over the input `result` and
// checks for.. well .. the upload status.
// it skips sessions that have been checked already.
func (u *Uploader) checkUploadSessionStatus(
	ctx context.Context,
	result []*sessionTypes.UploadSessionResult,
) (finished bool, err error) {
	var authedClient = u.ksefClient.GetAuthedHTTPClient()

	var processedSessions int
	for _, session := range result {
		if session.Processed {
			processedSessions++
			continue
		}

		sessionStatus, err := status.CheckSessionStatus(ctx, authedClient, session.SessionID)
		if err != nil {
			return true, err
		}
		session.Processed = sessionStatus.IsProcessed()
		if session.Processed {
			session.Status = sessionStatus
		}
	}
	return processedSessions == len(result), nil
}
