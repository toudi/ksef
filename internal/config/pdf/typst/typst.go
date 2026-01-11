package typst

type HeaderFooterConfig struct {
	Left   string `yaml:"left"`
	Center string `yaml:"center"`
	Right  string `yaml:"right"`
}

type TypstInvoicePrinterConfig struct {
	Template string             `yaml:"template"`
	Header   HeaderFooterConfig `yaml:"header"`
	Footer   HeaderFooterConfig `yaml:"footer"`
	Printout map[string]any     `yaml:"printout"`
}

type TypstUPOPrinterConfig struct {
	Template string `yaml:"template"`
}

type TypstPrinterConfig struct {
	Debug     bool                      `yaml:"debug,omitempty"`
	Workdir   string                    `yaml:"workdir,omitempty"`
	Templates string                    `yaml:"templates-dir"`
	Invoice   TypstInvoicePrinterConfig `yaml:"invoice,omitempty"`
	UPO       TypstUPOPrinterConfig     `yaml:"upo,omitempty"`
}
