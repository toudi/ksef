package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type StatusInfo struct {
	selectedFormat string
	sourcePath     string
	Environment    string              `json:"env" yaml:"env"`
	SessionID      string              `json:"sessionId" yaml:"sessionId"`
	InvoiceIds     []ksefInvoiceIdType `json:"invoiceIds,omitempty" yaml:"invoiceIds,omitempty"`
}

func StatusFromFile(filePath string) (*StatusInfo, error) {
	statusFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open status file: %v", err)
	}
	defer statusFile.Close()

	status := &StatusInfo{
		sourcePath:     filePath,
		selectedFormat: strings.TrimPrefix(strings.ToLower(filepath.Ext(filePath)), "."),
	}

	if status.selectedFormat == StatusFileFormatJSON {
		err = json.NewDecoder(statusFile).Decode(&status)
	} else if status.selectedFormat == StatusFileFormatYAML {
		err = yaml.NewDecoder(statusFile).Decode(&status)
	}

	return status, err
}

func (s *StatusInfo) Save(outputDirectory string) error {
	if s.selectedFormat == "" {
		s.selectedFormat = StatusFileFormatYAML
	}
	if outputDirectory == "" {
		if s.sourcePath == "" {
			return errors.New("please select output directory")
		}
		outputDirectory = filepath.Dir(s.sourcePath)
	}

	if s.selectedFormat != StatusFileFormatJSON && s.selectedFormat != StatusFileFormatYAML {
		return fmt.Errorf("unexpected format (expecting either `json` or `yaml`)")
	}

	outputFileName := filepath.Join(outputDirectory, fmt.Sprintf("status.%s", s.selectedFormat))
	statusFile, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("unable to create status file: %v", err)
	}
	defer statusFile.Close()

	if s.selectedFormat == StatusFileFormatJSON {
		return json.NewEncoder(statusFile).Encode(s)
	}
	if s.selectedFormat == StatusFileFormatYAML {
		return yaml.NewEncoder(statusFile).Encode(s)
	}

	return errors.New("unexpected return")
}
