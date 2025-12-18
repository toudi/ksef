package invoicesdb

import (
	"context"
	"errors"
	statuschecker "ksef/internal/invoicesdb/status-checker"
	statuscheckerconfig "ksef/internal/invoicesdb/status-checker/config"
)

var (
	errUnableToDiscoverPendingSessions = errors.New("unable to discover pending sessions")
)

func (i *InvoicesDB) checkPendingUploadSessions(
	ctx context.Context,
	statusCheckerConfig statuscheckerconfig.StatusCheckerConfig,
) error {
	// initialize checker:
	checker := statuschecker.NewStatusChecker(
		i.vip,
		i.ksefClient,
		statusCheckerConfig,
		i.monthsRange,
	)
	if err := checker.DiscoverPendingSessions(); err != nil {
		return errors.Join(errUnableToDiscoverPendingSessions, err)
	}

	return checker.CheckSessions(ctx)
}
