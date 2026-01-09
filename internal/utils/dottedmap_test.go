package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDottedMap(t *testing.T) {
	srcMap := map[string]string{
		"foo":         "bar",
		"baz.bzz.abc": "[1,2,3]",
		"baz.bzz.def": "84",
		"gbl":         "{\"a\": \"b\", \"c\": \"d\"}",
	}

	expectedMap := map[string]any{
		"foo": "bar",
		"baz": map[string]any{
			"bzz": map[string]any{
				"abc": []any{uint64(1), uint64(2), uint64(3)},
				"def": uint64(84),
			},
		},
		"gbl": map[string]any{"a": "b", "c": "d"},
	}

	dstMap, err := ReconstructMapFromDottedNotation(srcMap)
	require.NoError(t, err)
	require.Equal(t, expectedMap, dstMap)
}
