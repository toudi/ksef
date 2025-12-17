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

func SaveYAML(data any, dest string) error {
	writer, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer writer.Close()

	return yaml.NewEncoder(writer).Encode(data)
}
