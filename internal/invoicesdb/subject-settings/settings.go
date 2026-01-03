package subjectsettings

import (
	"errors"
	"ksef/internal/config"
	"ksef/internal/utils"
	"os"
	"path"
	"path/filepath"
)

const (
	settingsName = "settings.yaml"
)

var (
	errOpeningSettings         = errors.New("unable to open subject settings")
	errReadingSettingsContents = errors.New("unable to parse subject settings contents")
)

type SubjectSettings struct {
	PrinterSetup config.PDFPrinterConfig `yaml:"printer-config,omitempty"`
	JPK          *JPKSettings            `yaml:"jpk,omitempty"`
	dir          string
	dirty        bool
}

func OpenOrCreate(dir string) (*SubjectSettings, error) {
	settingsFile, exists, err := utils.FileExists(path.Join(dir, settingsName))
	if err != nil && !os.IsNotExist(err) {
		// the only way for the err to be not nil is when there's a problem opening
		// file
		return nil, errOpeningSettings
	}
	ss := &SubjectSettings{
		dir: dir,
		JPK: &JPKSettings{},
	}

	if exists {
		// if the file exists, then we need to read it's contents
		defer settingsFile.Close()
		if err = utils.ReadYAML(settingsFile, ss); err != nil {
			return nil, errors.Join(errReadingSettingsContents, err)
		}
	}

	return ss, nil
}

func (ss *SubjectSettings) Save() error {
	if !ss.dirty {
		return nil
	}

	return utils.SaveYAML(ss, filepath.Join(ss.dir, settingsName))
}
