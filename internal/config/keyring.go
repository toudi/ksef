package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

type KeyringEngine string

const (
	cfgKeyKeyringEngine             = "keyring.engine"
	cfgKeyKeyringFileLocation       = "keyring.file.path"
	cfgKeyKeyringFileBuffered       = "keyring.file.buffered"
	cfgKeyKeyringFileAskPassword    = "keyring.file.ask-password"
	cfgKeyKeyringFilePasswordFile   = "keyring.file.password-file"
	cfgKeyKeyringFilePasswordEnvVar = "keyring.file.password-env-var"
	keyringEngineSystem             = "system"
	keyringEngineFile               = "file"
)

var (
	errInvalidKeyringConfig               = errors.New("invalid keyring configuration")
	ErrEmptyPassword                      = errors.New("keyring password is empty")
	ErrEitherEnvVarOrPasswordFileRequired = errors.New("either env variable or password file have to be specified")
	ErrPasswordFilePermissionsTooWide     = errors.New("uprawnienia pliku hasła są zbyt szerokie. ustaw je na 0600")
)

type FileBasedKeyringConfig struct {
	Path     string // path to the keyring itself
	Buffered bool
	Password string // password to file
}

type Keyring struct {
	Engine KeyringEngine
	File   *FileBasedKeyringConfig
}

func init() {
	viper.SetDefault(cfgKeyKeyringEngine, keyringEngineSystem)
}

func FileKeyringFlags(flags *pflag.FlagSet) {
	flags.String(cfgKeyKeyringFileLocation, "", "ścieżka do keyringu opartego o plik")
	flags.Bool(cfgKeyKeyringFileAskPassword, false, "pytaj o hasło do keyringu na wejściu standardowym (stdin)")
	flags.Bool(cfgKeyKeyringFileBuffered, false, "buforuj keyring w pamięci")
	flags.String(cfgKeyKeyringFilePasswordFile, "", "ścieżka do pliku z hasłem keyringu")
	flags.String(cfgKeyKeyringFilePasswordEnvVar, "", "nazwa zmiennej środowiskowej która zawiera hasło do keyringu")
}

func KeyringFlags(flags *pflag.FlagSet) {
	flags.String(cfgKeyKeyringEngine, "system", "silnik keyringu")
	FileKeyringFlags(flags)
}

func GetFileBasedKeyringConfig(vip *viper.Viper) (*FileBasedKeyringConfig, error) {
	var cfg = &FileBasedKeyringConfig{
		Path:     viper.GetString(cfgKeyKeyringFileLocation),
		Buffered: viper.GetBool(cfgKeyKeyringFileBuffered),
	}
	if vip.GetBool(cfgKeyKeyringFileAskPassword) {
		// we have to ask user for the password
		fmt.Printf("podaj hasło do keyringu: \n")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil, err
		}
		cfg.Password = string(bytePassword)
	} else {
		envVarName := vip.GetString(cfgKeyKeyringFilePasswordEnvVar)
		passFile := vip.GetString(cfgKeyKeyringFilePasswordFile)

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

func KeyringConfig(vip *viper.Viper) (Keyring, error) {
	var err error
	var keyringConfig = Keyring{
		Engine: KeyringEngine(vip.GetString(cfgKeyKeyringEngine)),
	}

	if keyringConfig.Engine == keyringEngineFile {
		if keyringConfig.File, err = GetFileBasedKeyringConfig(vip); err != nil {
			return keyringConfig, err
		}
	}

	return keyringConfig, nil
}
