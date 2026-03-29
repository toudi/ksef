package keyring

import (
	"bytes"
	"errors"
	"ksef/internal/logging"
	"ksef/internal/utils"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

// file-based keyring implementation
type FileBasedKeyringConfig struct {
	Path     string // path to the keyring itself
	Buffered bool
	Password string // password to file
}

type FileBasedKeyring struct {
	cfg      *FileBasedKeyringConfig
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

func NewFileBasedKeyring(config *FileBasedKeyringConfig) (*FileBasedKeyring, error) {
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
	keyName := keyName(bucket, nip, key)
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
	keyName := keyName(bucket, nip, key)
	if f.cfg.Buffered {
		if value, exists := f.contents[keyName]; exists {
			return value, nil
		}
		return "", ErrNotFound
	}
	return f.retrieve(keyName)
}

func (f *FileBasedKeyring) Set(bucket string, nip string, key string, contents string) error {
	keyName := keyName(bucket, nip, key)
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
	plaintext, err := utils.GCMAESDecrypt(cipherText, []byte(f.cfg.Password))
	if err != nil {
		return nil, err
	}
	keyringContents := make(map[string]string)
	if err = yaml.NewDecoder(bytes.NewReader(plaintext)).Decode(&keyringContents); err != nil {
		return nil, errors.Join(errReading, err)
	}
	return keyringContents, nil
}

func (f *FileBasedKeyring) saveKeyring(contents map[string]string) error {
	var plaintextBuffer bytes.Buffer

	if err := yaml.NewEncoder(&plaintextBuffer).Encode(contents); err != nil {
		return err
	}

	ciphertext, err := utils.GCMAESEncrypt(plaintextBuffer.Bytes(), []byte(f.cfg.Password))
	if err != nil {
		logging.KeyringLogger.Error("error encrypting keyring contents", "err", err)
		return err
	}
	if err = os.WriteFile(f.cfg.Path, ciphertext, 0600); err != nil {
		logging.KeyringLogger.Error("error writing keyring to disk", "err", err)
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
