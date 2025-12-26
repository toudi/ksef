package commands

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"ksef/cmd/ksef/flags"
	"ksef/internal/certsdb"
	"ksef/internal/invoicesdb/config"
	"ksef/internal/keyring"
	"ksef/internal/logging"
	"ksef/internal/utils"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"time"

	encryptedZIP "github.com/alexmullins/zip"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var backupCommand = &cobra.Command{
	Use:     "backup",
	Short:   "zarchiwizuj katalogi stanowe",
	RunE:    backupRun,
	PreRunE: prepareBackupEnv,
}

var setBackupPasswordCommand = &cobra.Command{
	Use:   "set-password",
	Short: "ustaw hasło dla archiwów",
	RunE:  setBackupPassword,
}

var (
	errPasswordsDoNotMatch = errors.New("hasła różnią się")
	errOpeningKeyring      = errors.New("błąd inicjalizacji keyringu")
	errOpeningCertsDB      = errors.New("błąd otwierania bazy certyfikatów")
	errReadingPassword     = errors.New("błąd odczytu hasła archiwizacji")
	errWalkDir             = errors.New("błąd przeglądania katalogu")
	errCreatingArchive     = errors.New("błąd tworzenia pliku archiwum")
	errCreatingEntry       = errors.New("błąd dodawania pliku do archiwum")
	errCopyingToArchive    = errors.New("błąd kopiowania pliku do archiwum")
)

type BackupSourceType uint8

type BackupSource struct {
	sourceType BackupSourceType
	path       string
	certUIDs   []string
}

const (
	flagNameInvoicePDF                        = "invoice.pdf"
	flagNameUPO                               = "upo"
	flagNameUPOPDF                            = "upo.pdf"
	flagNameDate                              = "add-date"
	flagNameClear                             = "clear"
	backupSourceInvoices     BackupSourceType = 0x01
	backupSourceCertificates BackupSourceType = 0x02
	globalNIP                                 = "*"
)

type BackupConfig struct {
	InvoicePDF bool
	UPO        bool
	UPOPDF     bool
	sources    []BackupSource
	nip        string
	password   string
}

var (
	ring         keyring.Keyring
	excludeFiles = []string{".DS_Store"}
	backupConfig BackupConfig
)

func init() {
	backupFlags := backupCommand.PersistentFlags()
	backupLocalFlags := backupCommand.Flags()
	config.InvoicesDBFlags(backupLocalFlags)
	flags.NIP(backupFlags)
	backupLocalFlags.StringP(flagNameOutput, "o", "ksef-backup.zip", "plik wyjścia")
	backupLocalFlags.Bool(flagNameInvoicePDF, false, "archiwizuj pliki PDF faktur")
	backupLocalFlags.Bool(flagNameUPO, true, "archiwizuj UPO")
	backupLocalFlags.Bool(flagNameUPOPDF, false, "archiwizuj pliki PDF UPO")
	backupLocalFlags.BoolP(flagNameDate, "d", false, "dodaj datę do nazwy pliku backupu")
	backupLocalFlags.SortFlags = false
	passwordFlags := setBackupPasswordCommand.Flags()
	passwordFlags.BoolP(flagNameClear, "d", false, "usuń hasło")
	flags.NIP(passwordFlags)
	backupCommand.AddCommand(setBackupPasswordCommand)
}

func prepareBackupEnv(cmd *cobra.Command, _ []string) (err error) {
	vip := viper.GetViper()
	// for retrieving the password (optionally)
	ring, err = keyring.NewKeyring(vip)
	if err != nil {
		return errors.Join(errOpeningKeyring, err)
	}

	backupConfig = BackupConfig{
		InvoicePDF: vip.GetBool(flagNameInvoicePDF),
		UPO:        vip.GetBool(flagNameUPO),
		UPOPDF:     vip.GetBool(flagNameUPOPDF),
	}

	// check if the backup is supposed to be for a single NIP or for all of them
	backupConfig.nip = vip.GetString(flags.FlagNameNIP)
	invoicesDBConfig := config.GetInvoicesDBConfig(vip)
	backupConfig.sources = append(backupConfig.sources, BackupSource{
		sourceType: backupSourceInvoices, path: invoicesDBConfig.Root,
	})
	certsBackupSource := BackupSource{
		sourceType: backupSourceCertificates,
		path:       "certificates",
	}
	if backupConfig.nip != "" {
		certsDB, err := certsdb.OpenOrCreate(vip)
		if err != nil {
			return errors.Join(errOpeningCertsDB, err)
		}
		certsBackupSource.certUIDs = certsDB.FetchUIDsByNIP(backupConfig.nip)
		if backupConfig.password, err = ring.Get(keyring.AppPrefix, backupConfig.nip, keyring.KeyBackupPassword); err != nil && err != keyring.ErrNotFound {
			return errors.Join(errReadingPassword, err)
		}
	}
	backupConfig.sources = append(backupConfig.sources, certsBackupSource)
	// now, check if there's a password set
	if backupConfig.password == "" {
		if backupConfig.password, err = ring.Get(keyring.AppPrefix, globalNIP, keyring.KeyBackupPassword); err != nil && err != keyring.ErrNotFound {
			return errors.Join(errReadingPassword, err)
		}
	}
	if backupConfig.password == "" {
		logging.SeiLogger.Warn("brak zdefiniowanego hasła do archiwizacji. Archiwum nie będzie zaszyfrowane")
	}
	return nil
}

func backupRun(cmd *cobra.Command, _ []string) (err error) {
	logger := logging.BackupLogger.With("cmd", "backup")
	outputPath, err := cmd.Flags().GetString(flagNameOutput)
	if err != nil {
		return err
	}
	if viper.GetViper().GetBool(flagNameDate) {
		outputPath = strings.TrimSuffix(outputPath, ".zip")
		outputPath = outputPath + "-" + time.Now().Format(time.DateOnly)
		outputPath = outputPath + ".zip"
	}
	var fileWriter io.Writer
	var encryptedWriter *encryptedZIP.Writer
	var zipWriter *zip.Writer
	logger.Debug("tworzę plik archiwum", "output", outputPath)
	if fileWriter, err = os.Create(outputPath); err != nil {
		return errors.Join(errCreatingArchive, err)
	}
	if backupConfig.password != "" {
		logger.Debug("tryb szyfrowany")
		encryptedWriter = encryptedZIP.NewWriter(fileWriter)
	} else {
		zipWriter = zip.NewWriter(fileWriter)
	}

	configReader, exists, _ := utils.FileExists("config.yaml")
	if exists {
		var fileWriter io.Writer
		if encryptedWriter != nil {
			fileWriter, err = encryptedWriter.Encrypt("config.yaml", backupConfig.password)
			if err != nil {
				return errors.Join(errCreatingEntry, err)
			}
		} else {
			fileWriter, err = zipWriter.Create("config.yaml")
			if err != nil {
				return errors.Join(errCreatingEntry, err)
			}
		}
		if _, err = io.Copy(fileWriter, configReader); err != nil {
			return errors.Join(errCopyingToArchive, err)
		}
		configReader.Close()
	}

	for _, source := range backupConfig.sources {
		logger.Debug("przeglądam katalog", "path", source.path)
		if err = filepath.WalkDir(source.path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return errors.Join(errWalkDir, err)
			}

			if d.IsDir() || slices.Contains(excludeFiles, d.Name()) {
				// we're only interested in files
				return nil
			}

			extension := filepath.Ext(d.Name())
			isUPO := strings.Contains(path, "/upo/")

			if extension == ".pdf" {
				// either it's an UPO PDF but the used did not request it,
				// or it's an invoice PDF
				if isUPO && !backupConfig.UPOPDF || !backupConfig.InvoicePDF {
					return nil
				}
			}

			if isUPO && !backupConfig.UPO {
				// user does not want to backup UPO at all
				return nil
			}

			// check if we need to include the file
			if source.sourceType == backupSourceCertificates {
				if backupConfig.nip != "" {
					// if we only want to backup certificates for a specific nip and
					// this certificate is for somebody else, then simply return nil
					// to continue
					var skip bool = true

					for _, uid := range source.certUIDs {
						if strings.Contains(path, uid) {
							skip = false
							break
						}
					}

					if skip {
						return nil
					}
				}
			} else {
				if backupConfig.nip != "" && !strings.Contains(path, backupConfig.nip) {
					return nil
				}
			}

			var fileWriter io.Writer
			if encryptedWriter != nil {
				fileWriter, err = encryptedWriter.Encrypt(path, backupConfig.password)
				if err != nil {
					return errors.Join(errCreatingEntry, err)
				}
			} else {
				fileWriter, err = zipWriter.Create(path)
				if err != nil {
					return errors.Join(errCreatingEntry, err)
				}
			}
			if err = utils.CopyFileToWriter(path, fileWriter); err != nil {
				return errors.Join(errCopyingToArchive, err)
			}

			return nil
		}); err != nil {
			return errors.Join(errWalkDir, err)
		}
	}

	if encryptedWriter != nil {
		encryptedWriter.Flush()
		return encryptedWriter.Close()
	} else {
		zipWriter.Flush()
		return zipWriter.Close()
	}
}

func setBackupPassword(cmd *cobra.Command, _ []string) error {
	logger := logging.BackupLogger.With("cmd", "set-password")
	vip := viper.GetViper()
	nip := vip.GetString(flags.FlagNameNIP)
	clearPassword := vip.GetBool(flagNameClear)

	ring, err := keyring.NewKeyring(vip)
	if err != nil {
		return errors.Join(errOpeningKeyring, err)
	}

	var password string

	if !clearPassword {
		password, err = readUserPassword()
		if err != nil {
			return errors.Join(errReadingPassword, err)
		}
	}

	if nip != "" {
		if clearPassword {
			logger.Debug("usuwam hasło archiwizacji dla NIPu", "nip", nip)
			return ring.Delete(keyring.AppPrefix, nip, keyring.KeyBackupPassword)
		}
		logger.Debug("ustawiam hasło archiwizacji dla NIP", "nip", nip)
		return ring.Set(keyring.AppPrefix, nip, keyring.KeyBackupPassword, password)
	} else {
		if clearPassword {
			logger.Debug("usuwam globalne hasło archiwizacji")
			return ring.Delete(keyring.AppPrefix, globalNIP, keyring.KeyBackupPassword)
		}
		logger.Debug("ustawiam globalne hasło archiwizacji")
		return ring.Set(keyring.AppPrefix, globalNIP, keyring.KeyBackupPassword, password)
	}
}

func readUserPassword() (string, error) {
	fmt.Printf("podaj hasło do archiwizacji: \n")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Printf("powtórz hasło (celem weryfikacji): \n")
	bytePassword1, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	if !bytes.Equal(bytePassword, bytePassword1) {
		return "", errPasswordsDoNotMatch
	}
	return string(bytePassword), nil
}
