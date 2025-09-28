package archive

import "ksef/internal/utils"

type ArchivePart struct {
	utils.FilesizeAndHash
	FileName string
}
