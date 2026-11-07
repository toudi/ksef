package utils

import (
	"io"
	"os"

	"github.com/goccy/go-yaml"
)

func ReadYAML(reader io.Reader, dest any) error {
	err := yaml.NewDecoder(reader).Decode(dest)
	if err == io.EOF {
		return nil
	}
	return err
}

// WriteYAML encodes data as YAML into writer.
func WriteYAML(writer io.Writer, data any) error {
	return yaml.NewEncoder(writer).Encode(data)
}

func SaveYAML(data any, dest string) error {
	writer, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer writer.Close()

	return WriteYAML(writer, data)
}
