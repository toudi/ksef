package statuschecker

import (
	"errors"
	annualregistry "ksef/internal/invoicesdb/annual-registry"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	sessionregistry "ksef/internal/invoicesdb/session-registry"
	"ksef/internal/logging"
	"os"
	"slices"
)

var (
	errOpeningSessionRegistry = errors.New("unable to open session registry")
	errOpeningMonthlyRegistry = errors.New("unable to open monthly registry")
	errOpeningAnnualRegistry  = errors.New("unable to open annual registry")
)

func (c *StatusChecker) DiscoverPendingSessions() error {
	c.sessionIdToSessionRegistry = make(map[string]*sessionregistry.Registry)
	c.invoiceHashToMonthlyRegistry = make(map[string]*monthlyregistry.Registry)
	c.invoiceHashToAnnualRegistry = make(map[string]*annualregistry.Registry)

	yearToAnnualRegistryMap := make(map[int]*annualregistry.Registry)

	// here's the thing: we have our months range starting from previous month and ending with
	// current month. So what would happen is if we have invoices from previous month being
	// sent in the current month ? we would first iterate over the previous month, find the
	// checksums but because they would not be sent in previous month - we would not register
	// them in invoiceHashToMonthlyRegistry.
	// solution is quite simple - start from current month and end in previous month. this way,
	// we will make sure to visit all sessions and invoices.
	// fixes #25
	for _, month := range slices.Backward(c.monthsRange) {
		sessionRegistryForMonth, err := sessionregistry.OpenForMonth(c.vip, month)
		if err != nil {
			if os.IsNotExist(err) {
				logging.InvoicesDBLogger.Info("upload sessions registry does not exist for month", "month", month)
			} else {
				return errors.Join(errOpeningSessionRegistry, err)
			}
		} else {
			logging.InvoicesDBLogger.Debug("discovered upload sessions registry", "dir", sessionRegistryForMonth.Dir())
			for _, session := range sessionRegistryForMonth.PendingUploadSessions() {
				logging.InvoicesDBLogger.Debug("discovered upload session", "ref-no", session.RefNo)
				c.sessionIdToSessionRegistry[session.RefNo] = sessionRegistryForMonth
				// now that we have the session, we can iterate over invoice hashes so that
				// later we can bind them to monthly registries
				for _, invoice := range session.Invoices {
					if _, exists := c.invoiceHashToMonthlyRegistry[invoice.Checksum]; exists {
						logging.InvoicesDBLogger.Debug("invoice already bound to monthly registry and/or nil ptr")
						continue
					}
					// initialize to empty hash
					logging.InvoicesDBLogger.Debug("set monthly registry to nil ptr", "invoice", invoice.Checksum)
					c.invoiceHashToMonthlyRegistry[invoice.Checksum] = nil
					logging.InvoicesDBLogger.Debug("set annual registry to nil ptr", "invoice", invoice.Checksum)
					c.invoiceHashToAnnualRegistry[invoice.Checksum] = nil
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
			// check annual registry
			if _, exists := yearToAnnualRegistryMap[month.Year()]; !exists {
				annualRegistry, err := annualregistry.OpenForMonth(c.vip, month)
				if err != nil {
					if os.IsNotExist(err) {
						logging.InvoicesDBLogger.Info("annual registry does not exist for month", "month", month)
						continue
					}
					return errors.Join(errOpeningAnnualRegistry, err)
				}
				yearToAnnualRegistryMap[month.Year()] = annualRegistry
			}

			logging.InvoicesDBLogger.Debug("opened monthly registry", "path", monthlyRegistry.Dir())
			for invoiceHash := range c.invoiceHashToMonthlyRegistry {
				if monthlyRegistry.ContainsHash(invoiceHash) {
					logging.InvoicesDBLogger.Debug("override monthly registry for invoice", "invoice", invoiceHash)
					c.invoiceHashToMonthlyRegistry[invoiceHash] = monthlyRegistry
					c.invoiceHashToAnnualRegistry[invoiceHash] = yearToAnnualRegistryMap[month.Year()]
				}
			}
		}
	}

	// TODO: remove after done with debugging
	for invoiceHash, registry := range c.invoiceHashToMonthlyRegistry {
		registryPath := "nil"
		if registry != nil {
			registryPath = registry.Dir()
		}
		logging.InvoicesDBLogger.Debug("invoice to registry mapping", "invoice", invoiceHash, "registry", registryPath)
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

func (c *StatusChecker) SetInvoiceHashToAnnualRegistry(
	invoiceHashToAnnualRegistry map[string]*annualregistry.Registry,
) {
	c.invoiceHashToAnnualRegistry = invoiceHashToAnnualRegistry
}
