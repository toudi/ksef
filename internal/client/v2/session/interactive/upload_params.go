package interactive

import "time"

type UploadParams struct {
	ForceUpload bool
	Wait        time.Duration
}
