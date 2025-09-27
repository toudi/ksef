package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"
)

func Sha256FileToString(input string) (string, error) {
	file, err := os.Open(input)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hasher := sha256.New()
	if _, err = io.Copy(hasher, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func Sha256Base64(input []byte) (string, error) {
	var base64Encoder = base64.StdEncoding

	hash := sha256.New()
	if _, err := hash.Write(input); err != nil {
		return "", err
	}
	return base64Encoder.EncodeToString(hash.Sum(nil)), nil
}
