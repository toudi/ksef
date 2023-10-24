package inputprocessors

type csvConfig struct {
	Delimiter              string
	EncodingConversionFile string
}

type InputProcessorConfig struct {
	CSV       csvConfig
	Generator string
}
