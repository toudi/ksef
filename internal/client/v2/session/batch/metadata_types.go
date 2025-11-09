package batch

import (
	"ksef/internal/utils"
)

type BatchArchivePart struct {
	utils.FilesizeAndHash
	Filename string
}

type BatchMetadataInfo struct {
	Archive utils.FilesizeAndHash
	Parts   []BatchArchivePart
}
