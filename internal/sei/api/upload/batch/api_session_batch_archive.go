package batch

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"ksef/internal/encryption"
	"os"
	"path"
	"path/filepath"
)

type Archive struct {
	dir      string
	file     *os.File
	writer   *zip.Writer
	hash     []byte
	filesize int
}

type EncryptedArchive struct {
	cipher *encryption.Cipher
	size   int
	hash   []byte
}

type FileMeta struct {
	size     int
	hash     string
	filename string
}

func Archive_init(sourcePath string) (*Archive, error) {
	var err error
	archive := &Archive{dir: sourcePath}
	archive.file, err = os.Create(archive.filename())
	if err != nil {
		return nil, err
	}
	archive.writer = zip.NewWriter(archive.file)

	return archive, nil
}

func (a *Archive) filename() string {
	return filepath.Join(a.dir, "metadata.zip")
}

func (a *Archive) Close() {
	if a.writer != nil {
		a.writer.Close()
	}
	if a.file != nil {
		a.file.Close()
	}
}

func (a *Archive) addFile(fileName string) (*FileMeta, error) {
	var err error
	fileEntry, err := a.writer.Create(path.Base(fileName))
	if err != nil {
		return nil, err
	}
	fileContents, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	fileEntry.Write(fileContents)

	hash := sha256.New()
	hash.Write(fileContents)

	return &FileMeta{
		hash:     fmt.Sprintf("%x", hash.Sum(nil)),
		size:     len(fileContents),
		filename: fileName,
	}, nil
}

func (a *Archive) encryptedArchiveFileName() string {
	return a.filename() + ".aes"
}

func (a *Archive) encrypt(cipher *encryption.Cipher) (*EncryptedArchive, error) {
	var err error

	// read .zip file
	srcFileBytes, err := os.ReadFile(a.filename())
	hash := sha256.New()
	hash.Write(srcFileBytes)
	checksum := hash.Sum(nil)
	a.hash = make([]byte, len(checksum))
	a.filesize = len(srcFileBytes)
	copy(a.hash, checksum)
	if err != nil {
		return nil, fmt.Errorf("error reading archive file: %v", err)
	}
	encryptedBytes := cipher.Encrypt(srcFileBytes, true)
	if err != nil {
		return nil, fmt.Errorf("error encrypting archive file: %v", err)
	}

	dstFile, err := os.Create(a.encryptedArchiveFileName())
	if err != nil {
		return nil, fmt.Errorf("error creating encrypted file: %v", err)
	}
	_, err = io.Copy(dstFile, bytes.NewReader(encryptedBytes))
	if err != nil {
		return nil, fmt.Errorf("error writing encrypted file: %v", err)
	}

	hash = sha256.New()
	hash.Write(encryptedBytes)

	encryptedArchive := &EncryptedArchive{
		cipher: cipher,
		size:   len(encryptedBytes),
	}

	checksum = hash.Sum(nil)
	encryptedArchive.hash = make([]byte, len(checksum))

	copy(encryptedArchive.hash, checksum)

	return encryptedArchive, nil
}
