package keyring

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"ksef/internal/config"
	"ksef/internal/logging"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

// file-based keyring implementation
type FileBasedKeyring struct {
	cfg      *config.FileBasedKeyringConfig
	contents map[string]string
	dirty    bool
}

var (
	errCorruptedKeyringFile     = errors.New("keyring file corrupted")
	errUnableToInitializeCipher = errors.New("unable to initialize cipher")
	errDecryption               = errors.New("decryption failed")
	errReading                  = errors.New("reading decrypted keyring failed")
)

func keyName(bucket, nip, key string) string {
	return strings.Join([]string{bucket, nip, key}, "|")
}

func NewFileBasedKeyring(config *config.FileBasedKeyringConfig) (*FileBasedKeyring, error) {
	logging.KeyringLogger.Debug("file-based keyring initialization")
	var err error
	kr := &FileBasedKeyring{
		cfg: config,
	}
	if config.Buffered {
		kr.contents = make(map[string]string)
	}
	if _, err = os.Stat(config.Path); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	// keyring does not exist yet - that's not a problem.
	if os.IsNotExist(err) {
		return kr, nil
	}
	kr.contents, err = kr.loadKeyring()
	logging.KeyringLogger.Debug("file-based keyring initialized", "err", err)
	return kr, err
}

func (f *FileBasedKeyring) Delete(bucket string, nip string, key string) error {
	var keyName = keyName(bucket, nip, key)
	if f.cfg.Buffered {
		delete(f.contents, keyName)
		return nil
	} else {
		return f.upsert(func(keyring map[string]string) {
			delete(keyring, keyName)
		})
	}
}

func (f *FileBasedKeyring) Get(bucket string, nip string, key string) (string, error) {
	var keyName = keyName(bucket, nip, key)
	if f.cfg.Buffered {
		if value, exists := f.contents[keyName]; exists {
			return value, nil
		}
		return "", ErrNotFound
	}
	return f.retrieve(keyName)
}

func (f *FileBasedKeyring) Set(bucket string, nip string, key string, contents string) error {
	var keyName = keyName(bucket, nip, key)
	logging.KeyringLogger.Debug("set key", "keyName", keyName)
	if f.cfg.Buffered {
		f.contents[keyName] = contents
		f.dirty = true
		return nil
	}
	return f.upsert(func(tmp_value map[string]string) {
		tmp_value[keyName] = contents
	})
}

// this function is only going to be called from the CLI which explicitly sets the buffred flag
// therefore we can be certain that the inner contents map will be populated.
func (f *FileBasedKeyring) ForEach(callback func(bucket, nip, key, contents string) error) error {
	var err error
	for encodedKey, value := range f.contents {
		keyParts := strings.Split(encodedKey, "|")
		if len(keyParts) != 3 {
			logging.KeyringLogger.Warn("unexpected key format", "key", encodedKey)
			continue
		}
		if err = callback(keyParts[0], keyParts[1], keyParts[2], value); err != nil {
			return err
		}
	}

	return nil
}

// low level functions

func (f *FileBasedKeyring) upsert(updateFunc func(keyring map[string]string)) error {
	keyring, err := f.loadKeyring()
	if err != nil {
		return err
	}

	updateFunc(keyring)

	return f.saveKeyring(keyring)
}

func (f *FileBasedKeyring) retrieve(keyName string) (string, error) {
	keyring, err := f.loadKeyring()
	if err != nil {
		return "", err
	}

	if value, exists := keyring[keyName]; !exists {
		return "", ErrNotFound
	} else {
		return value, nil
	}
}

func (f *FileBasedKeyring) loadKeyring() (map[string]string, error) {
	var err error
	cipherText, err := os.ReadFile(f.cfg.Path)
	if err != nil {
		return nil, err
	}
	if len(cipherText) < 12 {
		return nil, errCorruptedKeyringFile
	}
	block, err := aes.NewCipher([]byte(f.cfg.Password))
	if err != nil {
		return nil, errors.Join(errUnableToInitializeCipher, err)
	}
	ciph, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Join(errUnableToInitializeCipher, err)
	}
	// first 12 bytes are the nonce
	plaintext, err := ciph.Open(nil, cipherText[:12], cipherText[12:], nil)
	if err != nil {
		return nil, errors.Join(errDecryption, err)
	}
	var keyringContents = make(map[string]string)
	if err = yaml.NewDecoder(bytes.NewReader(plaintext)).Decode(&keyringContents); err != nil {
		return nil, errors.Join(errReading, err)
	}
	return keyringContents, nil
}

func (f *FileBasedKeyring) saveKeyring(contents map[string]string) error {
	block, err := aes.NewCipher([]byte(f.cfg.Password))
	if err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	var plaintextBuffer bytes.Buffer
	if err = yaml.NewEncoder(&plaintextBuffer).Encode(contents); err != nil {
		return err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintextBuffer.Bytes(), nil)

	keyringFile, err := os.OpenFile(f.cfg.Path, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer keyringFile.Close()
	logging.KeyringLogger.Debug("write nonce to file")
	if _, err = keyringFile.Write(nonce); err != nil {
		return err
	}
	logging.KeyringLogger.Debug("write encrypted content to file")
	if _, err = keyringFile.Write(ciphertext); err != nil {
		return err
	}
	return nil
}

func (f *FileBasedKeyring) Close() error {
	logging.KeyringLogger.Debug("file-based keyring close()", "buffered", f.cfg.Buffered, "dirty", f.dirty)
	if !f.cfg.Buffered || !f.dirty {
		return nil
	}

	return f.saveKeyring(f.contents)
}
