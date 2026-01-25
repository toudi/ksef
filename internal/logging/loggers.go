package logging

import (
	"io"
	"log/slog"
)

const (
	defaultLevel = slog.LevelError
)

// these are the actual loggers that the program can reference
// initialize the default values of the loggers so that they
// become valid slog.Logger objects and so that we can use them
// without exploding the code
func defaultLogger() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(
			io.Discard,
			&slog.HandlerOptions{
				Level: defaultLevel,
			},
		),
	)
}

var (
	AuthLogger        *slog.Logger = defaultLogger()
	BackupLogger      *slog.Logger = defaultLogger().With("module", "backup")
	InvoicesDBLogger  *slog.Logger = defaultLogger().With("module", "invoicesDB")
	RegistryLogger    *slog.Logger = defaultLogger()
	CertsDBLogger     *slog.Logger = defaultLogger().With("module", "certsDB")
	KeyringLogger     *slog.Logger = defaultLogger()
	SeiLogger         *slog.Logger = defaultLogger()
	GenerateLogger    *slog.Logger = defaultLogger()
	UploadLogger      *slog.Logger = defaultLogger()
	HTTPLogger        *slog.Logger = defaultLogger()
	InteractiveLogger *slog.Logger = defaultLogger()
	BatchLogger       *slog.Logger = defaultLogger()
	DownloadLogger    *slog.Logger = defaultLogger()
	UPOLogger         *slog.Logger = defaultLogger()
	ParserLogger      *slog.Logger = defaultLogger()
	PDFRendererLogger *slog.Logger = defaultLogger()
	JPKLogger         *slog.Logger = defaultLogger()
)

var logLevels = map[string]slog.Level{}

func init() {
	// populate the helper map so that we can alter the loggers after config
	// is read.
	loggers = map[string]*slog.Logger{
		"main":         SeiLogger,
		"backup":       BackupLogger,
		"invoicesdb":   InvoicesDBLogger,
		"regitry":      RegistryLogger,
		"certsdb":      CertsDBLogger,
		"keyring":      KeyringLogger,
		"auth":         AuthLogger,
		"http":         HTTPLogger,
		"generate":     GenerateLogger,
		"upload":       UploadLogger,
		"interactive":  InteractiveLogger,
		"batch":        BatchLogger,
		"download":     DownloadLogger,
		"upo":          UPOLogger,
		"parser":       ParserLogger,
		"pdf-renderer": PDFRendererLogger,
		"jpk":          JPKLogger,
	}

	for loggerName := range loggers {
		logLevels[loggerName] = defaultLevel
	}
}
