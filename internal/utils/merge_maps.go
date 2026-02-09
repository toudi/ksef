package utils

import "errors"

func MergeMaps(dest map[string]any, src map[string]any) error {
	for k, v := range src {
		dstV, exists := dest[k]
		if !exists {
			dest[k] = v
		} else {
			if destMap, ok := dstV.(map[string]any); !ok {
				return errors.New("nested map expected (dst)")
			} else {
				if srcMap, ok := v.(map[string]any); !ok {
					return errors.New("nested map expected (src)")
				} else {
					return MergeMaps(destMap, srcMap)
				}
			}
		}
	}

	return nil
}
