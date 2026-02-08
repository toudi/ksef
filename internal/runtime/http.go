package runtime

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	cfgKeyHttpPollTime = "http.poll-wait-time"
	cfgKeyHttpRetries  = "http.max-retries"
	DefaultHttpRetries = 3
)

func HttpFlags(flagSet *pflag.FlagSet) {
	flagSet.Duration(cfgKeyHttpPollTime, 2*time.Second, "czas oczekiwania przy pollingu HTTP dla długotrwałych operacji")
	flagSet.Int(cfgKeyHttpRetries, 3, "maksymalna liczba powtórzeń requestu HTTP podczas napotkania błędu")
}

func HttpPollWaitTime(vip *viper.Viper) time.Duration {
	return vip.GetDuration(cfgKeyHttpPollTime)
}

func HttpMaxRetries(vip *viper.Viper) int {
	return vip.GetInt(cfgKeyHttpRetries)
}
