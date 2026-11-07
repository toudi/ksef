package utils

import (
	"io"
	"os"
	"path/filepath"
)

// WriteFileAtomic replaces dest with the content produced by write.
//
// The content is written to a temporary file in the destination directory, then
// flushed (Sync) and renamed over dest, so an interrupted write can never leave a
// truncated file behind. Missing parent directories are created.
//
// write must not keep the writer around after it returns.
func WriteFileAtomic(dest string, write func(w io.Writer) error) error {
	dir := filepath.Dir(dest)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(dir, "."+filepath.Base(dest)+"-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmpFile.Name()
	defer func() {
		tmpFile.Close()
		os.Remove(tmpName)
	}()

	if err = write(tmpFile); err != nil {
		return err
	}
	if err = tmpFile.Sync(); err != nil {
		return err
	}
	if err = tmpFile.Close(); err != nil {
		return err
	}

	return os.Rename(tmpName, dest)
}
