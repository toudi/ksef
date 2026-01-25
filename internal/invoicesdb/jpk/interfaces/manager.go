package interfaces

import (
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/invoicesdb/shared"
)

type JPKManager interface {
	ItemHasRule(invoice *monthlyregistry.Invoice, hash shared.ItemHash, rule func(shared.JPKItemRule) bool) bool
}
