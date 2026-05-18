package interfaces

import (
	"ksef/internal/invoicesdb/annotations"
	subjectsettings "ksef/internal/invoicesdb/subject-settings"
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
	annotations *annotations.Annotations,
	subjectSettings *subjectsettings.SubjectSettings,
	reportDate time.Time,
) JPKGenerator
