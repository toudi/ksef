package inputprocessors

type CSVFormat struct {
	source string
}

func CSVDecoder_Init(source string) InputProcessor {
	return &CSVFormat{source: source}
}

func (csv *CSVFormat) FeedLine() ([]string, error) {
	return []string{}, nil
}
