package utils

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
)

var (
	errUnableToCreateDir = errors.New("unable to create directory")
)

func CopyFile(srcPath string, destPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)

	return err
}

func SaveBufferToFile(buffer bytes.Buffer, destFileName string) error {
	if err := os.MkdirAll(filepath.Dir(destFileName), 0755); err != nil {
		return err
	}

	destFile, err := os.Create(destFileName)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, &buffer)
	return err
}
