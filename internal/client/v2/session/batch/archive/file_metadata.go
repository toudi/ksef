package archive

import (
	"ksef/internal/encryption"
	"ksef/internal/utils"
	"os"
)

type ArchivePart struct {
	utils.FilesizeAndHash
	FileName          string
	EncryptedFilename string
}

func (a *ArchivePart) Encrypt(cipher *encryption.Cipher) error {
	content, err := os.ReadFile(a.FileName)
	if err != nil {
		return err
	}
	encrypted := cipher.Encrypt(content, true)
	a.FileSize = int64(len(encrypted))
	a.Hash, err = utils.Sha256Base64(encrypted)
	if err != nil {
		return err
	}
	a.EncryptedFilename = a.FileName + ".aes"
	return os.WriteFile(a.EncryptedFilename, encrypted, 0644)
}
