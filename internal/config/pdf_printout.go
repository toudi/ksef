package config

import (
	"bytes"
	"errors"
	"fmt"
	"ksef/internal/config/pdf/cirfmf"
	"ksef/internal/config/pdf/typst"

	"github.com/goccy/go-yaml"
	"github.com/spf13/viper"
)

const (
	cfgKeyPdf string = "pdf"
)

var errEngineNotFound = errors.New("could not found engine for selected usage")

type PDFEngineConfig struct {
	UsageRaw     any                       `yaml:"usage"`
	TypstConfig  *typst.TypstPrinterConfig `yaml:"typst"`
	CIRFMFConfig *cirfmf.PrinterConfig     `yaml:"cirfmf"`
}

func (c PDFEngineConfig) Usage() []string {
	if usageSlice, ok := c.UsageRaw.([]any); ok {
		var usageStringSlice []string
		for _, usageStr := range usageSlice {
			usageStringSlice = append(usageStringSlice, usageStr.(string))
		}
		return usageStringSlice
	} else if usageStr, ok := c.UsageRaw.(string); ok {
		return []string{usageStr}
	}
	return []string{}
}

type PDFPrinterConfig struct {
	engines    []PDFEngineConfig
	usageIndex map[string]int
}

func (c *PDFPrinterConfig) GetEngine(usage string) (*PDFEngineConfig, error) {
	for _, usageAccessor := range []string{usage, "*"} {
		if index, exists := c.usageIndex[usageAccessor]; exists {
			return &c.engines[index], nil
		}
	}

	return nil, errEngineNotFound
}

func GetPDFPrinterConfig(vip *viper.Viper) (config PDFPrinterConfig, err error) {
	rawEngines := vip.Get(cfgKeyPdf)
	// let's use a dirty little trick here. instead of decoding the structs by hand let's
	// simply re-encode this raw slice of map[string]any as yaml to a temporary buffer
	// and then decode it from memory to a ready structs.
	var buffer bytes.Buffer
	if err = yaml.NewEncoder(&buffer).Encode(rawEngines); err != nil {
		fmt.Printf("unable to encode raw engines: %v", err)
		return config, err
	}
	var engines []PDFEngineConfig
	if err = yaml.NewDecoder(&buffer).Decode(&engines); err != nil {
		fmt.Printf("unable to decode raw engines back into slice: %v", err)
		return config, err
	}
	config.usageIndex = make(map[string]int)
	for i, engine := range engines {
		for _, usage := range engine.Usage() {
			config.usageIndex[usage] = i
		}
	}
	config.engines = engines

	return config, err
}
