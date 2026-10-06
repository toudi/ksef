package interfaces

import (
	"time"

	"ksef/internal/invoicesdb/annotations"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	subjectsettings "ksef/internal/invoicesdb/subject-settings"
	"ksef/internal/xml"

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
	annotations *annotations.Annotations,
	subjectSettings *subjectsettings.SubjectSettings,
	reportDate time.Time,
) JPKGenerator
