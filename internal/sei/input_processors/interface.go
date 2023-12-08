package inputprocessors

import "ksef/internal/sei/parser"

type InputProcessorConstructor func(string) InputProcessor

type InputProcessor interface {
	Process(sourceFile string, parser *parser.Parser) error
}
