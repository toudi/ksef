package runtime

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	cfgKeyHttpPollTime = "http.poll-wait-time"
)

func FlagHttpPollWaitTime(flagSet *pflag.FlagSet) {
	flagSet.Duration(cfgKeyHttpPollTime, 2*time.Second, "czas oczekiwania przy pollingu HTTP dla długotrwałych operacji")
}

func HttpPollWaitTime(vip *viper.Viper) time.Duration {
	return vip.GetDuration(cfgKeyHttpPollTime)
}
