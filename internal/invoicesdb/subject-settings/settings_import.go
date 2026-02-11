package subjectsettings

type AutoCorrectionConfig struct {
	Enabled bool   `yaml:"enabled"`
	Scheme  string `yaml:"scheme"`
}

type ImportSettings struct {
	AutoCorrection AutoCorrectionConfig `yaml:"auto-correction,omitempty"`
}
