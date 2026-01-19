package sessionregistry

import (
	"errors"
	"ksef/internal/client/v2/session/status"
	"ksef/internal/client/v2/upo"
	"ksef/internal/invoicesdb/config"
	"ksef/internal/logging"
	"ksef/internal/runtime"
	"ksef/internal/utils"
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/spf13/viper"
)

const (
	registryFilename string = "upload-sessions.yaml"
)

var (
	errOpeningRegistryFile     = errors.New("error opening upload sessions registry file")
	errReadingRegistryContents = errors.New("error reading upload sessions registry contents")
)

type Invoice struct {
	RefNo     string   `yaml:"ref-no,omitempty"`
	KSeFRefNo string   `yaml:"ksef-ref-no,omitempty"`
	Filename  string   `yaml:"filename,omitempty,omitzero"`
	Checksum  string   `yaml:"checksum"`
	Errors    []string `yaml:"errors,omitempty,omitzero"`
}

type UploadSession struct {
	Timestamp time.Time              `yaml:"timestamp"`
	RefNo     string                 `yaml:"ref-no"`
	Processed bool                   `yaml:"processed,omitempty"`
	Status    *status.StatusResponse `yaml:"status,omitempty"`
	Invoices  []*Invoice             `yaml:"invoices"`
	UPO       []upo.UPODownloadPage  `yaml:"upo,omitempty,omitzero"`
}

func (us *UploadSession) IsPending() bool {
	return us.Status == nil || us.Status.Status.Code < 200
}

type Registry struct {
	sessions []*UploadSession `yaml:"sessions"`

	dir    string
	dirty  bool
	logger *slog.Logger
}

func OpenOrCreate(dirName string) (*Registry, error) {
	// TODO:
	// this is a repeating pattern. I am more than convinced that I should do it in a generic way
	// (i.e. something like OpenOrCreate[Registry], restore func())
	// but I am under a time pressure
	regFile, exists, err := utils.FileExists(path.Join(dirName, registryFilename))
	if err != nil && !os.IsNotExist(err) {
		// the only way for the err to be not nil is when there's a problem opening
		// file
		return nil, errors.Join(errOpeningRegistryFile, err)
	}

	reg := &Registry{
		sessions: make([]*UploadSession, 0),
		dir:      dirName,
		logger:   logging.RegistryLogger.With("path", path.Join(dirName, registryFilename)),
	}

	if exists {
		// if the file exists, then we need to read it's contents
		defer regFile.Close()
		if err = utils.ReadYAML(regFile, &reg.sessions); err != nil {
			return nil, errors.Join(errReadingRegistryContents, err)
		}
	}

	return reg, nil
}

func OpenForMonth(vip *viper.Viper, month time.Time) (*Registry, error) {
	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return nil, err
	}

	environmentId := runtime.GetEnvironmentId(vip)

	invoicesDBConfig := config.GetInvoicesDBConfig(vip)

	// there is no active registry - let's try to create it.
	path := path.Join(
		invoicesDBConfig.Root,
		environmentId,
		nip,
		month.Format("2006"),
		month.Format("01"),
	)

	if _, exists, _ := utils.FileExists(path); !exists {
		return nil, os.ErrNotExist
	}

	return OpenOrCreate(path)
}

func (r *Registry) Dir() string {
	return r.dir
}
