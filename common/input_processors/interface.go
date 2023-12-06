package inputprocessors

import "ksef/common"

type InputProcessorConstructor func(string) InputProcessor

type InputProcessor interface {
	Process(sourceFile string, parser *common.Parser) error
}
