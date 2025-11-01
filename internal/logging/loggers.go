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

var AuthLogger *slog.Logger = defaultLogger()
var CertsDBLogger *slog.Logger = defaultLogger().With("module", "certsdb")
var SeiLogger *slog.Logger = defaultLogger()
var GenerateLogger *slog.Logger = defaultLogger()
var UploadLogger *slog.Logger = defaultLogger()
var HTTPLogger *slog.Logger = defaultLogger()
var InteractiveLogger *slog.Logger = defaultLogger()
var BatchLogger *slog.Logger = defaultLogger()
var DownloadLogger *slog.Logger = defaultLogger()
var UPOLogger *slog.Logger = defaultLogger()
var ParserLogger *slog.Logger = defaultLogger()
var PDFRendererLogger *slog.Logger = defaultLogger()

var logLevels = map[string]slog.Level{}

func init() {
	// populate the helper map so that we can alter the loggers after config
	// is read.
	loggers = map[string]*slog.Logger{
		"main":         SeiLogger,
		"certsdb":      CertsDBLogger,
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
	}

	for loggerName, _ := range loggers {
		logLevels[loggerName] = defaultLevel
	}
}
