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

func FileSizeAndSha256Hash(input string) (*FilesizeAndHash, error) {
	file, err := os.Open(input)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	hasher := sha256.New()
	var output = &FilesizeAndHash{}
	if output.FileSize, err = io.Copy(hasher, file); err != nil {
		return nil, err
	}
	output.Hash = hex.EncodeToString(hasher.Sum(nil))
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
