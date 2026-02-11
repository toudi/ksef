package annualregistry

import (
	"strconv"
	"strings"
	"time"
)

func (r *Registry) GenerateCorrectionNumber(scheme string, issueDate time.Time) string {
	// number of corrections within the current month
	var correctionsCountPerMonth int
	// number of overall corrections within the year.
	var correctionsCount int

	for _, invoice := range r.invoices {
		for _, correction := range invoice.Corrections {
			if correction.GenerationTime.Year() == issueDate.Year() && correction.GenerationTime.Month() == issueDate.Month() {
				correctionsCountPerMonth += 1
			}
			correctionsCount += 1
		}
	}

	return strings.NewReplacer(
		"{count}", strconv.Itoa(correctionsCount+1),
		"{count-month}", strconv.Itoa(correctionsCountPerMonth+1),
		"{year}", issueDate.Format("2006"),
		"{month}", issueDate.Format("01"),
	).Replace(scheme)
}
