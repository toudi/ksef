package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"
)

type FilesizeAndHash struct {
	FileSize int64
	Hash     string
}

func Sha256File(filename string) (int64, []byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return -1, nil, err
	}
	defer file.Close()
	hasher := sha256.New()
	fileSize, err := io.Copy(hasher, file)
	if err != nil {
		return -1, nil, err
	}
	return fileSize, hasher.Sum(nil), nil
}

func FileSizeAndSha256Hash(input string) (*FilesizeAndHash, error) {
	var err error
	var hashBytes []byte
	var output = &FilesizeAndHash{}
	if output.FileSize, hashBytes, err = Sha256File(input); err != nil {
		return nil, err
	}
	output.Hash = hex.EncodeToString(hashBytes)
	return output, nil
}

func Sha256Base64(input []byte) (string, error) {
	var base64Encoder = base64.StdEncoding

	hash := sha256.New()
	if _, err := hash.Write(input); err != nil {
		return "", err
	}
	return base64Encoder.EncodeToString(hash.Sum(nil)), nil
}
