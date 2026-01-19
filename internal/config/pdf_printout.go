package config

import (
	"bytes"
	"errors"
	invoicesdbconfig "ksef/internal/invoicesdb/config"
	subjectsettings "ksef/internal/invoicesdb/subject-settings"
	"ksef/internal/logging"
	pdfConfig "ksef/internal/pdf/config"
	"ksef/internal/runtime"
	"path/filepath"
	"slices"
	"text/template"

	"github.com/goccy/go-yaml"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

const (
	cfgKeyPdf string = "pdf"
)

var errEngineNotFound = errors.New("could not found engine for selected usage")

type PDFPrinterConfig struct {
	engines []pdfConfig.PDFEngineConfig
}

func (c *PDFPrinterConfig) GetEngines() []pdfConfig.PDFEngineConfig {
	return c.engines
}

func (c *PDFPrinterConfig) GetEngine(usageSelector pdfConfig.UsageSelector) (*pdfConfig.PDFEngineConfig, error) {
	// first, let's narrow down possible choices based on usage selector
	logging.PDFRendererLogger.Debug("trying to select matching engine", "usage selector", usageSelector)
	potentialEngines := lo.Filter(c.engines, func(e pdfConfig.PDFEngineConfig, _ int) bool {
		var matches bool

		matches = slices.Contains(e.Usage(), usageSelector.Usage) || slices.Contains(e.Usage(), "*")
		if !matches {
			return false
		}

		// if it matches, then we can also narrow it down based on participants:
		if len(usageSelector.Participants) > 0 && e.Condition != "" {
			logging.PDFRendererLogger.Debug("found condition - evaluating", "condition", e.Condition)
			tmpl, err := template.New("").Parse(e.Condition)
			if err != nil {
				logging.PDFRendererLogger.Error("unable to parse condition template", "err", err)
				return false
			}
			var buffer bytes.Buffer
			if err := tmpl.Execute(&buffer, usageSelector.Participants); err != nil {
				logging.PDFRendererLogger.Error("unable to render condition template", "err", err)
				return false
			}
			logging.PDFRendererLogger.Debug("condition evaluated to", "result", buffer.String())
			matches = buffer.String() == "true"
		}

		return matches
	})

	if len(potentialEngines) == 0 {
		return nil, errEngineNotFound
	}

	// now that we have a list of engines, we can sort them so that the most generic one is at the bottom:
	slices.SortFunc(potentialEngines, func(c1, c2 pdfConfig.PDFEngineConfig) int {
		// if either of the engines is global - move it to the bottom.
		if c1.UsageRaw == "*" {
			return 1
		}
		if c2.UsageRaw == "*" {
			return -1
		}
		// if one of the engines contains condition - make it preferrable
		if c1.Condition != "" && c2.Condition == "" {
			return -1
		}
		if c1.Condition == "" && c2.Condition != "" {
			return 1
		}
		return 0
	})

	return &potentialEngines[0], nil
}

func GetPDFPrinterConfig(vip *viper.Viper) (config PDFPrinterConfig, err error) {
	rawEngines := vip.Get(cfgKeyPdf)
	// let's use a dirty little trick here. instead of decoding the structs by hand let's
	// simply re-encode this raw slice of map[string]any as yaml to a temporary buffer
	// and then decode it from memory to a ready structs.
	var buffer bytes.Buffer
	if err = yaml.NewEncoder(&buffer).Encode(rawEngines); err != nil {
		return config, err
	}
	var engines []pdfConfig.PDFEngineConfig
	if err = yaml.NewDecoder(&buffer).Decode(&engines); err != nil {
		return config, err
	}
	// now let's check if we can override the engines with the ones from
	// subject settings
	logging.PDFRendererLogger.Debug("checking pdf config from subject settings")
	if subjectEngines, err := loadSubjectEngines(vip); len(subjectEngines) > 0 && err == nil {
		logging.PDFRendererLogger.Debug("overwrite pdf config with the one from subject settings")
		engines = subjectEngines
	} else {
		if err != nil {
			logging.PDFRendererLogger.Error("error loading subject settings", "err", err)
		}
	}

	config.engines = engines

	return config, err
}

func loadSubjectEngines(vip *viper.Viper) ([]pdfConfig.PDFEngineConfig, error) {
	cfg := invoicesdbconfig.GetInvoicesDBConfig(vip)
	nip, err := runtime.GetNIP(vip)
	if err != nil {
		return nil, err
	}

	ss, err := subjectsettings.OpenOrCreate(
		filepath.Join(
			cfg.Root, runtime.GetEnvironmentId(vip), nip,
		),
	)

	return ss.PDF, err
}
