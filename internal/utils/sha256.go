package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

func Sha256Base64(input []byte) (string, error) {
	var base64Encoder = base64.StdEncoding

	hash := sha256.New()
	if _, err := hash.Write(input); err != nil {
		return "", err
	}
	return base64Encoder.EncodeToString(hash.Sum(nil)), nil
}
