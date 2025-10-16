package utils

import (
	"encoding/base64"
	"encoding/hex"
)

func Base64ToHex(input string) string {
	bytes, _ := base64.StdEncoding.DecodeString(input)
	return hex.EncodeToString(bytes)
}
