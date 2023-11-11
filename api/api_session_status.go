package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type StatusInfoFile struct {
	Environment string `json:"env" yaml:"env"`
	ReferenceNo string `json:"ref" yaml:"ref"`
}

func saveStatusInfo(info StatusInfoFile, outputDirectory string, format string) error {
	if format == StatusFileFormatJSON {
		statusFile, err := os.Create(filepath.Join(outputDirectory, "status.json"))
		if err != nil {
			return fmt.Errorf("unable to create status file: %v", err)
		}
		defer statusFile.Close()
		return json.NewEncoder(statusFile).Encode(info)
	} else if format == StatusFileFormatYAML {
		statusFile, err := os.Create(filepath.Join(outputDirectory, "status.yaml"))
		if err != nil {
			return fmt.Errorf("unable to create status file: %v", err)
		}
		defer statusFile.Close()
		return yaml.NewEncoder(statusFile).Encode(info)
	}

	return fmt.Errorf("unexpected format (expecting either `json` or `yaml`)")
}
