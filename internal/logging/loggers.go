package logging

import (
	"io"
	"log/slog"
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
				Level: slog.LevelInfo,
			},
		),
	)
}

var SeiLogger *slog.Logger = defaultLogger()
var GenerateLogger *slog.Logger = defaultLogger()
var UploadLogger *slog.Logger = defaultLogger()
var UploadHTTPLogger *slog.Logger = defaultLogger()
var InteractiveLogger *slog.Logger = defaultLogger()
var InteractiveHTTPLogger *slog.Logger = defaultLogger()
var BatchLogger *slog.Logger = defaultLogger()
var BatchHTTPLogger *slog.Logger = defaultLogger()
var DownloadHTTPLogger *slog.Logger = defaultLogger()
var DownloadLogger *slog.Logger = defaultLogger()
var UPOLogger *slog.Logger = defaultLogger()
var UPOHTTPLogger *slog.Logger = defaultLogger()
var ParserLogger *slog.Logger = defaultLogger()

func init() {
	// populate the helper map so that we can alter the loggers after config
	// is read.
	loggers = map[string]*slog.Logger{
		"main":             SeiLogger,
		"generate":         GenerateLogger,
		"upload":           UploadLogger,
		"upload.http":      UploadHTTPLogger,
		"interactive":      InteractiveLogger,
		"interactive.http": InteractiveHTTPLogger,
		"batch":            BatchLogger,
		"batch.http":       BatchHTTPLogger,
		"download":         DownloadLogger,
		"download.http":    DownloadHTTPLogger,
		"upo":              UPOLogger,
		"upo.http":         UPOHTTPLogger,
		"parser":           ParserLogger,
	}
}
