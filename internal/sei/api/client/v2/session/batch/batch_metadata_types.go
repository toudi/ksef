package batch

import "ksef/internal/sei/api/client/v2/session/batch/archive"

type BatchArchivePart struct {
	archive.FileSizeAndHash
	Filename string
}

type BatchMetadataInfo struct {
	Archive archive.FileSizeAndHash
	Parts   []BatchArchivePart
}
