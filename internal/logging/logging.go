package logging

import (
	"errors"
	"fmt"
	"io"
	"ksef/internal/config"
	"log/slog"
	"os"
)

var SeiLogger *slog.Logger
var GenerateLogger *slog.Logger

var loggers = map[string]*slog.Logger{
	"main":     SeiLogger,
	"generate": GenerateLogger,
}

var errUnknownLogger = errors.New("Unknown logger")
var outputWriter io.Writer
var outputFile *os.File

func parseLevel(logLevel string) slog.Level {
	switch logLevel {
	case "info":
		return slog.LevelInfo
	case "debug":
		return slog.LevelDebug
	default:
		return slog.LevelError
	}
}

func InitLogging(output string) error {
	if output == "" {
		return nil
	}

	config := config.Config
	if config.Logging != nil {
		var err error

		if output == "-" {
			outputWriter = os.Stdout
		} else {
			outputFile, err = os.OpenFile(output, os.O_CREATE, 0644)
			if err != nil {
				return fmt.Errorf("unable to open log file: %v", err)
			}
			outputWriter = outputFile
		}

		var logger *slog.Logger
		var exists bool

		for loggerName, logLevel := range config.Logging {
			if _, exists = loggers[loggerName]; exists {

				logger = slog.New(slog.NewTextHandler(outputWriter, &slog.HandlerOptions{
					Level: parseLevel(logLevel),
				}))
				loggers[loggerName] = logger

				// TODO: this is actually very smelly.
				//       what I was hoping for was to automatically change the
				//       reference of the original logger since it's a pointer
				//       but maybe we need to use unsafe.Pointer for that ?
				switch loggerName {
				case "main":
					SeiLogger = logger
				case "generate":
					GenerateLogger = logger
				default:
				}
			} else {
				return errUnknownLogger
			}
		}
	}

	return nil
}

func init() {
	SeiLogger = slog.Default()
	GenerateLogger = slog.Default()
}

func FinishLogging() {
	if outputFile != nil {
		outputFile.Close()
	}
}
