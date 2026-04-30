package generators

import (
	"fmt"
	"ksef/internal/interfaces"
	"ksef/internal/sei/generators/fa_3_1"
)

var registry map[string]interfaces.Generator

func registerGenerator(name string, g interfaces.Generator) {
	if registry == nil {
		registry = make(map[string]interfaces.Generator, 0)
	}

	registry[name] = g
}

func Generator(id string) (interfaces.Generator, error) {
	generator, exists := registry[id]
	if !exists {
		return nil, fmt.Errorf("unknown generator: %s", id)
	}

	return generator, nil
}

func init() {
	registerGenerator("fa-3_1.0", fa_3_1.GeneratorFactory())
}
