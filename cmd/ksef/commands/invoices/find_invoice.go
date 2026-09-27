package invoices

import (
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"time"

	"github.com/spf13/viper"
)

func findInvoiceByRefNo(vip *viper.Viper, refNo string) (*monthlyregistry.Invoice, *monthlyregistry.Registry, error) {
	month := time.Now()
	month = time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.Local)

	var invoice *monthlyregistry.Invoice
	var reg *monthlyregistry.Registry
	var err error

	// let's start from the most recent month and iterate back
	// until we hit a match.
	for range 12 {
		reg, err = monthlyregistry.OpenForMonth(vip, month)
		if err != nil {
			month = month.AddDate(0, -1, 0)
			continue
		}
		if invoice = reg.GetInvoice(func(i monthlyregistry.Invoice) bool {
			return i.RefNo == refNo && i.Type == monthlyregistry.InvoiceTypeIssued
		}); invoice != nil {
			break
		}

		month = month.AddDate(0, -1, 0)
	}

	if invoice == nil {
		return nil, nil, errInvoiceNotFound
	}

	return invoice, reg, nil
}
