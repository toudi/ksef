package utils_test

import (
	"ksef/internal/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimeParsingFromString(t *testing.T) {
	type unitTest struct {
		input    string
		expected time.Time
	}

	for _, test := range []unitTest{
		{
			input:    "2026-02-05",
			expected: time.Date(2026, 2, 5, 0, 0, 0, 0, time.Local),
		},
		{
			input:    "2026-02-05T13:14:15",
			expected: time.Date(2026, 2, 5, 13, 14, 15, 0, time.Local),
		},
		{
			input:    "2026-02-05 13:14:15",
			expected: time.Date(2026, 2, 5, 13, 14, 15, 0, time.Local),
		},
		{
			input:    "2026-02-05T13:14:15Z",
			expected: time.Date(2026, 2, 5, 13, 14, 15, 0, time.UTC),
		},
	} {
		t.Run(test.input, func(t *testing.T) {
			parsed, err := utils.ParseTimeFromString(test.input)
			require.NoError(t, err)
			require.Equal(t, test.expected, parsed)
		})
	}
}
