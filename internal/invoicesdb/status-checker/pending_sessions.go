package statuschecker

import (
	"errors"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	sessionregistry "ksef/internal/invoicesdb/session-registry"
	"ksef/internal/logging"
	"os"
)

var (
	errOpeningSessionRegistry = errors.New("unable to open session registry")
	errOpeningMonthlyRegistry = errors.New("unable to open monthly registry")
)

func (c *StatusChecker) DiscoverPendingSessions() error {
	c.sessionIdToSessionRegistry = make(map[string]*sessionregistry.Registry)
	c.invoiceHashToMonthlyRegistry = make(map[string]*monthlyregistry.Registry)

	for _, month := range c.monthsRange {
		sessionRegistryForMonth, err := sessionregistry.OpenForMonth(c.vip, month)
		if err != nil {
			if os.IsNotExist(err) {
				logging.InvoicesDBLogger.Info("upload sessions registry does not exist for month", "month", month)
			} else {
				return errors.Join(errOpeningSessionRegistry, err)
			}
		} else {
			for _, session := range sessionRegistryForMonth.PendingUploadSessions() {
				c.sessionIdToSessionRegistry[session.RefNo] = sessionRegistryForMonth
				// now that we have the session, we can iterate over invoice hashes so that
				// later we can bind them to monthly registries
				for _, invoice := range session.Invoices {
					// initialize to empty hash
					c.invoiceHashToMonthlyRegistry[invoice.Checksum] = nil
				}
			}
			// now we can open the monthly registry for this month and
			// check if it contains invoice hashes that are of interest for us.
			monthlyRegistry, err := monthlyregistry.OpenForMonth(c.vip, month)
			if err != nil {
				if os.IsNotExist(err) {
					logging.InvoicesDBLogger.Info("monthly registry does not exist for month", "month", month)
					continue
				}
				return errors.Join(errOpeningMonthlyRegistry, err)
			}
			for invoiceHash := range c.invoiceHashToMonthlyRegistry {
				if monthlyRegistry.ContainsHash(invoiceHash) {
					c.invoiceHashToMonthlyRegistry[invoiceHash] = monthlyRegistry
				}
			}
		}
	}

	return nil
}

func (c *StatusChecker) AddSessionID(sessionID string, registry *sessionregistry.Registry) {
	c.sessionIdToSessionRegistry[sessionID] = registry
}

func (c *StatusChecker) SetInvoiceHashToMonthlyRegistry(
	invoiceHashToMonthlyRegistry map[string]*monthlyregistry.Registry,
) {
	c.invoiceHashToMonthlyRegistry = invoiceHashToMonthlyRegistry
}
