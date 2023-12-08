package inputprocessors

type csvConfig struct {
	Delimiter              string
	EncodingConversionFile string
}

type xlsxConfig struct {
	SheetName string
}

type InputProcessorConfig struct {
	CSV       csvConfig
	XLSX      xlsxConfig
	Generator string
}
