package utils_test

import (
	"ksef/internal/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBase64ChunkedString(t *testing.T) {
	type testCase struct {
		input    []byte
		width    int
		expected string
	}

	for _, test := range []testCase{
		{
			input:    []byte{'a', 'a', 'a'},
			width:    96,
			expected: "YWFh",
		},
		{
			input:    []byte{'a', 'a', 'a'},
			width:    2,
			expected: "YW\nFh\n",
		},
		{
			input:    []byte{'a', 'a', 'a'},
			width:    3,
			expected: "YWF\nh\n",
		},
	} {
		t.Run(test.expected, func(t *testing.T) {
			t.Parallel()

			result := utils.Base64ChunkedString(test.input, test.width)
			require.Equal(t, test.expected, result)
		})
	}
}
