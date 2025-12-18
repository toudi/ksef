package archive

import (
	"archive/zip"
	"ksef/internal/utils"
	"os"
	"path/filepath"
)

type Archive struct {
	output     *os.File
	outputPath string
	writer     *zip.Writer
	// current size of the archive
	size int64
	// maximum allowed archive size
	maxFileSize int64
	// where the target zip archive should be created in
	outputDir string
	// base name of the file, without .zip extension
	// this will serve as a template for creating archive parts
	basename   string
	Parts      []ArchivePart
	partWriter *os.File
}

func New(basename string, maxFileSize int64) (*Archive, error) {
	outputFilename := basename + ".zip"
	output, err := os.Create(outputFilename)
	if err != nil {
		return nil, err
	}
	writer := zip.NewWriter(output)

	return &Archive{
		maxFileSize: maxFileSize,
		output:      output,
		outputPath:  outputFilename,
		outputDir:   filepath.Dir(outputFilename),
		writer:      writer,
	}, nil
}

func (a *Archive) Metadata() (*utils.FilesizeAndHash, error) {
	return utils.FileSizeAndSha256Hash(a.outputPath)
}
