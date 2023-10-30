package metadata

import (
	"ksef/common/aes"
)

type EncryptedArchive struct {
	cipher *aes.Cipher
	size   int
	hash   []byte
}
