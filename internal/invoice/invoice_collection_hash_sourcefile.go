package invoice

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

func hashSourceFile(fileName string) (string, error) {
	var err error

	hasher := sha256.New()
	content, err := os.ReadFile(fileName)

	if err != nil {
		return "", fmt.Errorf("unable to read file: %v", err)
	}

	hasher.Write(content)

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
