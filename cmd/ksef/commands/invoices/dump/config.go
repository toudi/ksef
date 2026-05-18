package dump

import (
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	cfgKeyTemplatePath = "monthly-dump.typst-template-path"
	cfgKeyGenerator    = "monthly-dump.generator"
)

const (
	DefaultGenerator    = "WSI Pegasus"
	DefaultTemplatePath = "examples/local-pdf-printout/typst/annotation/accountant-notes.typ"
)

// AnnotationConfig holds configuration for the accountant-notes Typst template.
type AnnotationConfig struct {
	TemplatePath string
	YamlPath     string
	Generator    string
}

// MonthlyDumpFlags registers flags for monthly dump configuration.
func MonthlyDumpFlags(flagSet *pflag.FlagSet) {
	flagSet.String(cfgKeyTemplatePath, DefaultTemplatePath, "ścieżka do szablonu Typst do generowania PDF z adnotacjami dla księgowych")
	flagSet.String(cfgKeyGenerator, DefaultGenerator, "nazwa generatora raportu")
}

// GetAnnotationConfig reads the annotation template configuration from viper.
func GetConfig(vip *viper.Viper) AnnotationConfig {
	config := AnnotationConfig{
		TemplatePath: DefaultTemplatePath,
		Generator:    DefaultGenerator,
	}
	if templatePath := vip.GetString(cfgKeyTemplatePath); templatePath != "" {
		config.TemplatePath = templatePath
	}
	if generator := vip.GetString(cfgKeyGenerator); generator != "" {
		config.Generator = generator
	}

	// YAML file is generated at runtime, always placed alongside the template
	config.YamlPath = filepath.Join(filepath.Dir(config.TemplatePath), "annotations.yaml")

	return config
}
