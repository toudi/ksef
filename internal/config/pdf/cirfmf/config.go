package cirfmf

type PrinterConfig struct {
	TemplatesDir string `yaml:"templates-dir"`
	WorkDir      string `yaml:"workdir"`
	NodeBinPath  string `yaml:"node-bin"`
}
