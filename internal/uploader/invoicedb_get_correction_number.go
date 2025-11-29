package uploader

import (
	"fmt"
	"ksef/internal/config"
	"strings"
	"time"
)

func (idb *InvoiceDB) getCorrectionNumber(issued time.Time) string {
	var now = time.Now()
	cfg := config.CorrectionsConfig(idb.vip)
	template := cfg.Numbering
	// let's prepare some placeholder values
	var correctionsInCurrentMonth int = 0
	var correctionsCount int = 0

	for _, invoice := range idb.Invoices {
		if !invoice.Correction {
			continue
		}
		correctionsCount += 1
		if invoice.GenerationTime.Month() == now.Month() {
			correctionsInCurrentMonth += 1
		}
	}

	return strings.NewReplacer(
		"{count}", fmt.Sprintf("%d", correctionsCount+1),
		"{currMonthCount}", fmt.Sprintf("%d", correctionsInCurrentMonth+1),
		"{year}", fmt.Sprintf("%d", now.Year()),
	).Replace(template)
}
