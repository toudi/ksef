package certsdb

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"
	kr "ksef/internal/keyring"
	"ksef/internal/utils"
	"os"
)

var encryptionHeader = []byte{
	'E', 'P', 'K', 0x00, 0x99,
}

// loadPrivateKey reads bytes of the underlying private key and
// decrypts them, if required.
func (c *Certificate) readPrivateKeyBytes(keyring kr.Keyring) ([]byte, error) {
	rawBytes, err := os.ReadFile(c.PrivateKeyFilename())
	if err != nil {
		return nil, err
	}

	if bytes.HasPrefix(rawBytes, encryptionHeader) {
		decryptionKey, err := c.loadEncryptionKey(keyring, c.UID)
		if err != nil {
			return nil, err
		}
		rawBytes, err = utils.GCMAESDecrypt(rawBytes[len(encryptionHeader):], decryptionKey)
	}

	return rawBytes, err
}

func (c *Certificate) loadEncryptionKey(keyring kr.Keyring, certID string) ([]byte, error) {
	// retrieve decryption key from keyring
	decryptionKey, err := keyring.Get(kr.AppPrefix, "", kr.PrivateKeyEncryptionKey(certID))
	if err != nil {
		return nil, err
	}
	// decryption key is base64 encoded since we're dealing with bytes but zalando lib stores strings so let's just be
	// safe than sorry
	return base64.StdEncoding.DecodeString(decryptionKey)
}

func (c *Certificate) PrivateKeyIsEncrypted() (bool, error) {
	privateKeyBytes, err := os.ReadFile(c.PrivateKeyFilename())
	if err != nil {
		return false, err
	}

	return bytes.HasPrefix(privateKeyBytes, encryptionHeader), nil
}

func (c *Certificate) encryptPrivateKeyBytes(pemBytes []byte, dest io.Writer, keyring kr.Keyring) error {
	encryptionKey := make([]byte, 16)
	_, err := rand.Read(encryptionKey)
	if err != nil {
		return err
	}
	// let's save the key to keyring first before proceeding further
	if err = keyring.Set(kr.AppPrefix, "", kr.PrivateKeyEncryptionKey(c.UID), base64.StdEncoding.EncodeToString(encryptionKey)); err != nil {
		return err
	}
	encryptedKeyBytes, err := utils.GCMAESEncrypt(pemBytes, encryptionKey)
	if err != nil {
		return err
	}
	encryptedKeyBytes = bytes.Join([][]byte{encryptionHeader, encryptedKeyBytes}, nil)
	_, err = io.Copy(dest, bytes.NewReader(pemBytes))
	return err
}

func (c *Certificate) EncryptPrivateKey(keyring kr.Keyring) error {
	privateKeyBytes, err := os.ReadFile(c.PrivateKeyFilename())
	if err != nil {
		return err
	}
	privateKeyFile, err := os.OpenFile(c.PrivateKeyFilename(), os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()
	return c.encryptPrivateKeyBytes(privateKeyBytes, privateKeyFile, keyring)
}
