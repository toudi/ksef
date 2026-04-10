package interfaces

import (
	"ksef/internal/invoicesdb/jpk/manager"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/xml"
	"time"

	"github.com/beevik/etree"
)

type JPKGenerator interface {
	Document() (*xml.Node, error)
	ProcessInvoice(
		invoice *monthlyregistry.Invoice,
		doc *etree.Document,
	) error
}

type JPKGeneratorFactory func(
	manager *manager.JPKManager,
	reportDate time.Time,
) JPKGenerator
