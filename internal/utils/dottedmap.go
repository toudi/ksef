package utils

import (
	"bytes"
	"strings"

	"github.com/goccy/go-yaml"
)

func ReconstructMapFromDottedNotation(source map[string]string) (dest map[string]any, err error) {
	dest = make(map[string]any)
	var buffer bytes.Buffer

	for key, value := range source {
		// find entry in dest map by the key
		var tmpDest any = dest
		var keyName string

		keyParts := strings.Split(key, ".")
		for idx, keyPart := range keyParts {
			_, exists := tmpDest.(map[string]any)[keyPart]
			if !exists {
				if idx < len(keyParts)-1 {
					tmpDest.(map[string]any)[keyPart] = make(map[string]any)
				}
			}
			keyName = keyPart
			if idx < len(keyParts)-1 {
				tmpDest = tmpDest.(map[string]any)[keyPart]
			}
		}

		// treat the value as a stringified yaml and try to decode
		buffer.Reset()
		buffer.WriteString(value)
		var dstValue any
		if err = yaml.NewDecoder(&buffer).Decode(&dstValue); err != nil {
			break
		}
		tmpDest.(map[string]any)[keyName] = dstValue
	}

	return dest, err
}
