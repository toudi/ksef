package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var (
	errUnableToDecrypt          = errors.New("unable to decrypt ciphertext")
	errUnableToInitializeCipher = errors.New("unable to initialize cipher")
	errDecryption               = errors.New("error during decryption")
)

// helper methods for dealing with gcm encryption/decryption
func GCMAESEncrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return aesgcm.Seal(nil, nonce, plaintext, nil), nil
}

func GCMAESDecrypt(ciphertext, key []byte) ([]byte, error) {
	if len(ciphertext) < 12 {
		return nil, errUnableToDecrypt
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Join(errUnableToInitializeCipher, err)
	}
	ciph, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Join(errUnableToInitializeCipher, err)
	}
	// first 12 bytes are the nonce
	plaintext, err := ciph.Open(nil, ciphertext[:12], ciphertext[12:], nil)
	if err != nil {
		return nil, errors.Join(errDecryption, err)
	}
	return plaintext, nil
}
