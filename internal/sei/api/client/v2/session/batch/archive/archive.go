package archive

import (
	"archive/zip"
	"os"
)

type Archive struct {
	output     *os.File
	outputPath string
	writer     *zip.Writer
	// current size of the archive
	size int
	// maximum allowed archive size
	maxFileSize int
	// where the target zip archive should be created in
	outputDir string
	// base name of the file, without .zip extension
	// this will serve as a template for creating archive parts
	basename   string
	Parts      []ArchivePart
	partWriter *os.File
}

func New(basename string, maxFileSize int) (*Archive, error) {
	output, err := os.Create(basename + ".zip")
	if err != nil {
		return nil, err
	}
	writer := zip.NewWriter(output)

	return &Archive{
		maxFileSize: maxFileSize,
		output:      output,
		outputPath:  basename + ".zip",
		writer:      writer,
	}, nil
}

func (a *Archive) Metadata() (FileSizeAndHash, error) {
	return FileSizeAndHash{}, nil
}
