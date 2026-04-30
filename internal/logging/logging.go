package logging

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

// this is a utility map that will be used when config file will be read
// and log level can be applied. Unfortunetely there's no way to change the
// log level and/or output at runtime therefore we have to re-initialize the
// logger.
var loggers = map[string]*slog.Logger{}

var (
	errUnknownLogger = errors.New("Unknown logger")
	outputWriter     io.Writer
	outputFile       *os.File
	Verbose          bool = false
)

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

func InitLogging(output string, vip *viper.Viper) error {
	if output == "" {
		return nil
	}

	var err error
	loggingConfig := LoggingConfig(vip)

	if output == "-" {
		outputWriter = os.Stdout
	} else {
		var logFileFlag int = os.O_APPEND
		if vip.GetBool(CfgKeyLogFileTruncate) {
			logFileFlag = os.O_TRUNC
		}
		outputFile, err = os.OpenFile(output, os.O_CREATE|os.O_RDWR|logFileFlag, 0644)
		if err != nil {
			return fmt.Errorf("unable to open log file: %v", err)
		}
		outputWriter = outputFile
	}

	// initiate loggers
	var logger *slog.Logger
	var loggerName string

	for loggerName, logger = range loggers {
		logLevel := logLevels[loggerName]

		if Verbose {
			logLevel = slog.LevelDebug
		} else {
			// let's see if the logger level was overriden via config:
			if level, exists := loggingConfig[loggerName]; exists {
				logLevel = parseLevel(level)
			} else if level, exists := loggingConfig["*"]; exists {
				logLevel = parseLevel(level)
			}
		}

		// it may look cryptic and ugly but the bottom line here is this:
		// we take `logger` which is a pointer to `slog.Logger` and we want to
		// re-initialize it, however we also want the address to stay the same.
		*logger = *slog.New(slog.NewTextHandler(outputWriter, &slog.HandlerOptions{
			Level: logLevel,
		}))

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
