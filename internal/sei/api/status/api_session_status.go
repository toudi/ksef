package status

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func StatusFromFile(filePath string) (*StatusInfo, error) {
	statusFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open status file: %v", err)
	}
	defer statusFile.Close()

	status := &StatusInfo{
		SourcePath:     filePath,
		SelectedFormat: strings.TrimPrefix(strings.ToLower(filepath.Ext(filePath)), "."),
	}

	if status.SelectedFormat == StatusFileFormatJSON {
		err = json.NewDecoder(statusFile).Decode(&status)
	} else if status.SelectedFormat == StatusFileFormatYAML {
		err = yaml.NewDecoder(statusFile).Decode(&status)
	}

	return status, err
}

func (s *StatusInfo) Save(outputDirectory string) error {
	if s.SelectedFormat == "" {
		s.SelectedFormat = StatusFileFormatYAML
	}
	if outputDirectory == "" {
		if s.SourcePath == "" {
			return errors.New("please select output directory")
		}
		outputDirectory = filepath.Dir(s.SourcePath)
	}

	if s.SelectedFormat != StatusFileFormatJSON && s.SelectedFormat != StatusFileFormatYAML {
		return fmt.Errorf("unexpected format (expecting either `json` or `yaml`)")
	}

	outputFileName := filepath.Join(outputDirectory, fmt.Sprintf("status.%s", s.SelectedFormat))
	statusFile, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("unable to create status file: %v", err)
	}
	defer statusFile.Close()

	if s.SelectedFormat == StatusFileFormatJSON {
		return json.NewEncoder(statusFile).Encode(s)
	}
	if s.SelectedFormat == StatusFileFormatYAML {
		return yaml.NewEncoder(statusFile).Encode(s)
	}

	return errors.New("unexpected return")
}
