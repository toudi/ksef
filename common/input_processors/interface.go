package inputprocessors

type InputProcessorConstructor func(string) InputProcessor

type InputProcessor interface {
	FeedLine() ([]string, error)
}
