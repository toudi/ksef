package logging

import (
	"errors"
	"fmt"
	"io"
	"ksef/internal/config"
	"log/slog"
	"os"
)

// this is a utility map that will be used when config file will be read
// and log level can be applied. Unfortunetely there's no way to change the
// log level and/or output at runtime therefore we have to re-initialize the
// logger.
var loggers = map[string]*slog.Logger{}

var errUnknownLogger = errors.New("Unknown logger")
var outputWriter io.Writer
var outputFile *os.File
var Verbose bool = false

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

	var err error

	if output == "-" {
		outputWriter = os.Stdout
	} else {
		outputFile, err = os.OpenFile(output, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return fmt.Errorf("unable to open log file: %v", err)
		}
		outputWriter = outputFile
	}

	if Verbose {
		for _, loggerPtr := range loggers {
			var logger *slog.Logger = loggerPtr
			*logger = *slog.New(slog.NewTextHandler(outputWriter, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))
		}
		return nil
	}

	config := config.GetConfig()

	if config.Logging != nil {
		var logger *slog.Logger
		var exists bool

		for loggerName, logLevel := range config.Logging {
			if logger, exists = loggers[loggerName]; exists {
				// it may look cryptic and ugly but the bottom line here is this:
				// we take `logger` which is a pointer to `slog.Logger` and we want to
				// re-initialize it, however we also want the address to stay the same.
				*logger = *slog.New(slog.NewTextHandler(outputWriter, &slog.HandlerOptions{
					Level: parseLevel(logLevel),
				}))
			} else {
				return errUnknownLogger
			}
		}
	}

	return nil
}

func FinishLogging() {
	if outputFile != nil {
		if err := outputFile.Close(); err != nil {
			fmt.Printf("error cosing logfile: %v\n", err)
		}
	}
}
