package fa

import (
	"ksef/internal/sei/constants"
	"ksef/internal/xml"
	"strings"
	"time"
)

type ElementOrdering map[string]map[string]int

type GeneratorHooks struct {
	PostProcess func(root *xml.Node) error
}

var defaultHooks = GeneratorHooks{
	PostProcess: func(root *xml.Node) error { return nil },
}

// naiwnie łudzę się, że ministerstwo nie będzie drastycznie zmieniać struktury dokumentu
// dzięki czemu nie będę musiał w nieskończoność kopiować generatora a jedyne co się zmieni
// to wartości nagłówków
type FAGenerator struct {
	commonData map[string]string
	// niestety, ministerstwo używa typu sequence co oznacza, że musimy odpowiednio posortować
	// elementy drzewa - inaczej dokument nie przejdzie walidacji
	elementOrdering ElementOrdering
	hooks           GeneratorHooks
	runTimestamp    time.Time
}

func New(initializers ...func(fa *FAGenerator)) *FAGenerator {
	generator := &FAGenerator{
		runTimestamp: time.Now().UTC(),
		hooks:        defaultHooks,
	}

	for _, init := range initializers {
		init(generator)
	}

	return generator
}

func (fg *FAGenerator) IssuerTIN() string {
	return fg.commonData["Faktura.Podmiot1.DaneIdentyfikacyjne.NIP"]
}

func (fg *FAGenerator) isCommonData(section string) bool {
	sectionLower := strings.ToLower(section)

	return (sectionLower == constants.SectionInvoiceRoot ||
		sectionLower == constants.SectionInvoiceHeader ||
		strings.HasPrefix(sectionLower, constants.SectionInvoiceIssuer))
}

func (fg *FAGenerator) isItemSection(section string) bool {
	return strings.ToLower(section) == constants.SectionInvoiceItemRow
}
