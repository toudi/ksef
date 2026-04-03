package config

import (
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestDownloaderConfig(t *testing.T) {
	t.Run("make sure it does not fail when endDate != nil", func(t *testing.T) {
		t.Parallel()

		vip := viper.New()
		vip.Set(flagEndDate, "2026-02-05")
		vip.Set(flagStartDate, "2026-02-04")

		expectedEndDate := time.Date(2026, 2, 5, 0, 0, 0, 0, time.Local)

		params, err := GetDownloaderConfig(vip, "")
		require.NoError(t, err)

		require.NotNil(t, params.EndDate)
		require.Equal(t, expectedEndDate, *params.EndDate)
	})

	t.Run("allow time as input as well", func(t *testing.T) {
		t.Parallel()

		vip := viper.New()
		vip.Set(flagStartDate, "2026-02-05")
		vip.Set(flagEndDate, "2026-02-06T13:14:15Z")

		expectedEndDate := time.Date(2026, 2, 6, 13, 14, 15, 0, time.UTC)
		params, err := GetDownloaderConfig(vip, "")
		require.NoError(t, err)

		require.NotNil(t, params.EndDate)
		require.Equal(t, expectedEndDate, *params.EndDate)
	})

	t.Run("do not allow endDate set before startDate", func(t *testing.T) {
		t.Parallel()

		vip := viper.New()
		vip.Set(flagStartDate, "2026-02-05")
		vip.Set(flagEndDate, "2026-02-01")

		_, err := GetDownloaderConfig(vip, "")
		require.Equal(t, errEndDateBeforeStartDate, err)
	})
}
