package keyring

import (
	"errors"
	"fmt"
	"ksef/internal/flags"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/viper"
	"golang.org/x/term"
)

var (
	errInvalidKeyringConfig               = errors.New("invalid keyring configuration")
	ErrEmptyPassword                      = errors.New("keyring password is empty")
	ErrEitherEnvVarOrPasswordFileRequired = errors.New("either env variable or password file have to be specified")
	ErrPasswordFilePermissionsTooWide     = errors.New("uprawnienia pliku hasła są zbyt szerokie. ustaw je na 0600")
)

func NewKeyring(vip *viper.Viper) (Keyring, error) {
	engine := vip.GetString(flags.CfgKeyKeyringEngine)

	if engine == flags.KeyringEngineFile {
		cfg, err := GetFileBasedKeyringConfig(vip)
		if err != nil {
			return nil, err
		}
		return NewFileBasedKeyring(cfg)
	}

	return NewSystemKeyring(), nil
}

func GetFileBasedKeyringConfig(vip *viper.Viper) (*FileBasedKeyringConfig, error) {
	cfg := &FileBasedKeyringConfig{
		Path:     vip.GetString(flags.CfgKeyKeyringFileLocation),
		Buffered: vip.GetBool(flags.CfgKeyKeyringFileBuffered),
	}
	if vip.GetBool(flags.CfgKeyKeyringFileAskPassword) {
		// we have to ask user for the password
		fmt.Printf("podaj hasło do keyringu: \n")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil, err
		}
		cfg.Password = string(bytePassword)
	} else {
		envVarName := vip.GetString(flags.CfgKeyKeyringFilePasswordEnvVar)
		passFile := vip.GetString(flags.CfgKeyKeyringFilePasswordFile)

		if envVarName == "" && passFile == "" {
			return nil, errors.Join(errInvalidKeyringConfig, ErrEitherEnvVarOrPasswordFileRequired)
		}

		if envVarName != "" {
			// password will be provided via env variable
			cfg.Password = os.Getenv(envVarName)
		} else {
			// password will be provided via file.
			// let's make sure that it has permissions set to 0600 - otherwise let's bail out
			stat, err := os.Stat(passFile)
			if err != nil {
				return nil, err
			}
			if stat.Mode() != 0600 {
				return nil, errors.Join(errInvalidKeyringConfig, ErrPasswordFilePermissionsTooWide)
			}
			passwordBytes, err := os.ReadFile(passFile)
			if err != nil {
				return nil, err
			}
			cfg.Password = string(passwordBytes)
		}
	}

	if cfg.Password == "" {
		return nil, errors.Join(errInvalidKeyringConfig, ErrEmptyPassword)
	}
	cfg.Password = strings.TrimSpace(cfg.Password)
	return cfg, nil
}
