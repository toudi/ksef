package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type FilepathResolverConfig struct {
	Path            string
	Mkdir           bool
	DefaultFilename string
}

var ErrDoesNotExistAndMkdirNotSpecified = errors.New("")

func ResolveFilepath(config FilepathResolverConfig) (string, error) {
	// if, by any chance, the caller entered a valid filename which just does not exist
	// yet, let's assign it to the final output value.
	var finalOutput string = config.Path
	var err error
	// let's validate output.
	// first, let's check if this is a file or a directory.
	outputExt := filepath.Ext(config.Path)
	outputPath := filepath.Dir(config.Path)

	if outputExt == "" {
		// since there is no filename extension we have to treat the whole thing as a path.
		outputPath = config.Path
		finalOutput = filepath.Join(
			outputPath,
			config.DefaultFilename,
		)
	}

	// let's validate output directory
	_, err = os.Stat(outputPath)

	if os.IsNotExist(err) {
		// that's still fine at this point. let's check if we can create it.
		if !config.Mkdir {
			return "", ErrDoesNotExistAndMkdirNotSpecified
		}
		if err = os.MkdirAll(outputPath, 0755); err != nil {
			return "", fmt.Errorf("unable to create directories: %v", err)
		}
	}

	return finalOutput, nil
}
