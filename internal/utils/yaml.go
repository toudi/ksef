package utils

import (
	"io"
	"os"

	"github.com/goccy/go-yaml"
)

func ReadYAML(reader io.Reader, dest any) error {
	return yaml.NewDecoder(reader).Decode(dest)
}

func SaveYAML(data any, dest string) error {
	writer, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer writer.Close()

	return yaml.NewEncoder(writer).Encode(data)
}
