package inputprocessors

import (
	"errors"
	"io"
	"ksef/internal/sei/parser"
)

var errNotImplemented = errors.New("not implemented")

type InputProcessorConstructor func(string) InputProcessor

type InputProcessor interface {
	Process(sourceFile string, parser *parser.Parser) error
	ProcessReader(src io.Reader, parser *parser.Parser) error
}
