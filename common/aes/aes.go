package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

type Cipher struct {
	Key []byte
	IV  []byte

	block cipher.Block
}

func CipherInit(keysize int) (*Cipher, error) {
	var err error

	_cipher := &Cipher{}
	_cipher.Key = make([]byte, keysize)
	_cipher.IV = make([]byte, aes.BlockSize)

	_, err = rand.Read(_cipher.Key)
	if err != nil {
		return nil, fmt.Errorf("could not initialize encryption cipher: %v", err)
	}

	if _, err = rand.Read(_cipher.IV); err != nil {
		return nil, fmt.Errorf("could not initialize cipher's IV: %v", err)
	}

	_cipher.block, err = aes.NewCipher(_cipher.Key)
	if err != nil {
		return nil, fmt.Errorf("could not initialize cipher block: %v", err)
	}

	return _cipher, nil
}

func (c *Cipher) Encrypt(input []byte, add_pkcs7pad bool) []byte {
	plaintext := input
	if add_pkcs7pad {
		plaintext, _ = pkcs7Pad(input, c.block.BlockSize())
	}

	ciphertext := make([]byte, len(plaintext))

	bm := cipher.NewCBCEncrypter(c.block, c.IV)
	bm.CryptBlocks(ciphertext, plaintext)

	return ciphertext

}
