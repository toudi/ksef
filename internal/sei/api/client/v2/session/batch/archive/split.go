package archive

import (
	"fmt"
	"io"
	"os"
)

// https://github.com/CIRFMF/ksef-docs/blob/main/sesja-wsadowa.md#2-podzia%C5%82-binarny-paczki-zip-na-cz%C4%99%C5%9Bci
// as I understand the tutorial, it's not really splitting but rather chunking
// maybe it is a limitation of the zip package itself? no idea.
func (a *Archive) Split(maxPartSize int) error {
	// first, let's sure that the file is written as a whole:
	var err error

	a.writer.Close()
	var reader *os.File
	if reader, err = os.Open(a.outputPath); err != nil {
		return err
	}
	defer reader.Close()

	var chunkWriter *HasherWriter
	if chunkWriter, err = NewHasherWriter(reader, maxPartSize, a.getPartWriter, a.chunkWritten); err != nil {
		return err
	}

	_, err = io.Copy(chunkWriter, reader)
	if err == nil {
		_, _ = chunkWriter.Write(nil)
	}

	return nil
}

func (a *Archive) getPartWriter() (io.Writer, error) {
	var err error

	if a.partWriter != nil {
		a.partWriter.Close()
	}

	a.Parts = append(a.Parts, ArchivePart{
		FileName: fmt.Sprintf("%s.zip.%03d", a.basename, len(a.Parts)),
	})
	a.partWriter, err = os.Create(a.Parts[len(a.Parts)-1].FileName)

	return a.partWriter, err
}

func (a *Archive) chunkWritten(hash string, bytesWritten int) {
	a.Parts[len(a.Parts)-1].FileSize = uint64(bytesWritten)
	a.Parts[len(a.Parts)-1].Hash = hash
}
