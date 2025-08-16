package generators

import (
	"fmt"
	"ksef/internal/interfaces"
	"ksef/internal/sei/generators/fa_2"
	"ksef/internal/sei/generators/fa_3_1"
)

var registry map[string]interfaces.Generator

func registerGenerator(name string, g interfaces.Generator) {
	if registry == nil {
		registry = make(map[string]interfaces.Generator, 0)
	}

	registry[name] = g
}

// func Run(generatorName string, delimiter string, inputFile string, outputDirectory string, encodingConversionFile string) (interfaces.Generator, error) {
// 	var g interfaces.Generator
// 	g, exists := registry[generatorName]
// 	if !exists {
// 		return nil, fmt.Errorf("unknown generator")
// 	}

// 	var input *os.File
// 	var err error

// 	parser := &ksef.Parser{LineHandler: g.LineHandler, Comma: delimiter, EncodingConversionFile: encodingConversionFile}
// 	input, err = os.Open(inputFile)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot open input file")
// 	}
// 	defer input.Close()

// 	if err = parser.Parse(input); err != nil {
// 		return nil, fmt.Errorf("unable to parse file: %v", err)
// 	}

// 	if _, err = os.Stat(outputDirectory); os.IsNotExist(err) {
// 		if err = os.MkdirAll(outputDirectory, 0755); err != nil {
// 			return nil, fmt.Errorf("cannot create output directory: %v", err)
// 		}
// 	}

// 	return g, g.Save(outputDirectory)
// }

func Generator(id string) (interfaces.Generator, error) {
	generator, exists := registry[id]
	if !exists {
		return nil, fmt.Errorf("unknown generator: %s", id)
	}

	return generator, nil
}

func init() {
	registerGenerator("fa-2", fa_2.GeneratorFactory())
	registerGenerator("fa-3_1.0", fa_3_1.GeneratorFactory())
}
