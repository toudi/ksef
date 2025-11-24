package utils

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
)

func Base64ToHex(input string) string {
	bytes, _ := base64.StdEncoding.DecodeString(input)
	return hex.EncodeToString(bytes)
}

func Base64ChunkedString(input []byte, chunkLength int) string {
	var inputBase64 = base64.StdEncoding.EncodeToString(input)
	if len(inputBase64) < chunkLength {
		return inputBase64
	}
	var output strings.Builder

	for len(inputBase64) > 0 {
		if len(inputBase64) < chunkLength {
			chunkLength = len(inputBase64)
		}
		output.WriteString(inputBase64[:chunkLength])
		output.WriteString("\n")
		inputBase64 = inputBase64[chunkLength:]
	}

	return output.String()
}
