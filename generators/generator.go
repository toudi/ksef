package generators

import (
	"fmt"
	"ksef"
	"os"
)

type Generator interface {
	LineHandler(string, map[string]string) error
	Save(string) error
	Issuer() string
}

var registry map[string]Generator

func registerGenerator(name string, g Generator) error {
	if registry == nil {
		registry = make(map[string]Generator, 0)
	}

	registry[name] = g

	return nil
}

func Run(generatorName string, delimiter string, inputFile string, outputDirectory string) (Generator, error) {
	var g Generator
	g, exists := registry[generatorName]
	if !exists {
		return nil, fmt.Errorf("unknown generator")
	}

	var input *os.File
	var err error

	parser := &ksef.Parser{LineHandler: g.LineHandler, Comma: delimiter}
	input, err = os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open input file")
	}
	defer input.Close()

	if err = parser.Parse(input); err != nil {
		return nil, fmt.Errorf("unable to parse file: %v", err)
	}

	if _, err = os.Stat(outputDirectory); os.IsNotExist(err) {
		if err = os.MkdirAll(outputDirectory, 0755); err != nil {
			return nil, fmt.Errorf("cannot create output directory: %v", err)
		}
	}

	return g, g.Save(outputDirectory)
}
