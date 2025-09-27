package archive

type FileSizeAndHash struct {
	FileSize uint64
	Hash     string
}

type ArchivePart struct {
	FileSizeAndHash
	FileName string
}
