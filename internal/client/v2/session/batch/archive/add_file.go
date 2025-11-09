package archive

import (
	"errors"
	"io"
	"os"
	"path"
)

var (
	ErrExeedsMaxSize = errors.New("archive exceeds maximum filesize")
)

func (a *Archive) AddFile(fileName string) error {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return err
	}

	if (a.size + int(fileInfo.Size())) > a.maxFileSize {
		return ErrExeedsMaxSize
	}

	if err = a.addFileToArchive(fileName); err != nil {
		return err
	}

	return nil
}

func (a *Archive) addFileToArchive(fileName string) error {
	fileEntry, err := a.writer.Create(path.Base(fileName))
	if err != nil {
		return err
	}
	// copy the file
	fileObj, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fileObj.Close()
	if _, err = io.Copy(fileEntry, fileObj); err != nil {
		return err
	}
	return nil
}
