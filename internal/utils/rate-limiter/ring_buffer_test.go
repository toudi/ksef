package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCallHistoryPush(t *testing.T) {
	t.Run("test writeIndex behaviour", func(t *testing.T) {
		b := NewCallHistory(4)
		b.Push(time.Date(2026, 2, 1, 2, 3, 4, 5, time.Local))
		b.Push(time.Date(2026, 2, 1, 2, 3, 4, 6, time.Local))
		b.Push(time.Date(2026, 2, 1, 2, 3, 4, 7, time.Local))
		require.Equal(t, 3, b.writeIndex)
		b.Push(time.Date(2026, 2, 1, 2, 3, 4, 8, time.Local))
		require.Equal(t, 0, b.writeIndex)
		require.Equal(t, []time.Time{
			time.Date(2026, 2, 1, 2, 3, 4, 5, time.Local),
			time.Date(2026, 2, 1, 2, 3, 4, 6, time.Local),
			time.Date(2026, 2, 1, 2, 3, 4, 7, time.Local),
			time.Date(2026, 2, 1, 2, 3, 4, 8, time.Local),
		}, b.entries)
	})
}

func TestReplacingLastEntry(t *testing.T) {
	t.Run("single entry", func(t *testing.T) {
		b := NewCallHistory(1)
		b.Push(time.Date(2026, 2, 8, 23, 54, 0, 0, time.Local))
		expected := time.Date(2026, 2, 8, 23, 54, 2, 0, time.Local)
		b.ReplaceLastEntry(expected)
		require.Equal(t, []time.Time{expected}, b.entries)
		require.Equal(t, b.writeIndex, 0)
	})

	t.Run("two entries", func(t *testing.T) {
		b := NewCallHistory(2)
		firstValue := time.Date(2026, 2, 8, 23, 54, 0, 0, time.Local)
		b.Push(firstValue)
		b.Push(time.Date(2026, 2, 8, 23, 54, 2, 0, time.Local))
		expected := time.Date(2026, 2, 8, 23, 54, 4, 0, time.Local)
		b.ReplaceLastEntry(expected)
		require.Equal(t, []time.Time{firstValue, expected}, b.entries)
		require.Equal(t, b.writeIndex, 0)
	})
}
