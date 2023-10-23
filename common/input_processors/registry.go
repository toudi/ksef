package inputprocessors

var processorsRegistry = map[string]InputProcessorConstructor{
	"csv": CSVDecoder_Init,
	// "xlsx": XLSXDecoder_Init,
}
