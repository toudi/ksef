package utils_test

import (
	"ksef/internal/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergeMaps(t *testing.T) {
	map1 := map[string]any{
		"printout": map[string]any{
			"nested": map[string]any{
				"foo": "bar",
			},
			"test": "a-key",
		},
	}
	map2 := map[string]any{
		"printout": map[string]any{
			"nested": map[string]any{
				"def": "abc",
			},
		},
	}

	require.NoError(t, utils.MergeMaps(map1, map2))

	require.Equal(t, map[string]any{
		"printout": map[string]any{
			"nested": map[string]any{
				"foo": "bar",
				"def": "abc",
			},
			"test": "a-key",
		},
	}, map1)
}
