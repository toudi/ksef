package statuschecker

import (
	v2 "ksef/internal/client/v2"
	"ksef/internal/client/v2/session/status"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	sessionregistry "ksef/internal/invoicesdb/session-registry"
	"ksef/internal/invoicesdb/status-checker/config"
	"time"

	"github.com/spf13/viper"
)

type StatusChecker struct {
	vip         *viper.Viper
	checker     *status.SessionStatusChecker
	cfg         config.StatusCheckerConfig
	monthsRange []time.Time
	// internal helper maps
	sessionIdToSessionRegistry   map[string]*sessionregistry.Registry
	invoiceHashToMonthlyRegistry map[string]*monthlyregistry.Registry
}

func NewStatusChecker(
	vip *viper.Viper,
	ksefClient *v2.APIClient,
	cfg config.StatusCheckerConfig,
	monthsRange []time.Time,
) *StatusChecker {
	return &StatusChecker{
		vip:         vip,
		checker:     ksefClient.SessionStatusChecker(),
		cfg:         cfg,
		monthsRange: monthsRange,
	}
}
