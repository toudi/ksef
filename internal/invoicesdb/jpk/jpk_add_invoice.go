package jpk

import (
	"fmt"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
)

func (j *JPK) AddIncome(invoice *monthlyregistry.XMLInvoice) error {
	fmt.Printf("invoice: %+v\n", invoice)
	return nil
}

func (j *JPK) AddReceived(invoice *monthlyregistry.XMLInvoice) error {
	return nil
}

func (j *JPK) Save(output string) error {
	return nil
}
