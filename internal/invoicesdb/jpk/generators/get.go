package generators

import (
	"ksef/internal/invoicesdb/jpk/generators/interfaces"
	"ksef/internal/invoicesdb/jpk/generators/jpk_v7m_3"
	"ksef/internal/invoicesdb/jpk/manager"
	"time"
)

type generator struct {
	factory        interfaces.JPKGeneratorFactory
	availableSince time.Time
	id             string
}

var availableGenerators = []generator{
	{
		id:             "v7m_3",
		factory:        jpk_v7m_3.New,
		availableSince: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
	},
}

func GetJPKGenerator(
	manager *manager.JPKManager,
	reportMonth time.Time,
) interfaces.JPKGenerator {
	var selected interfaces.JPKGeneratorFactory
	today := time.Now().Local()

	for _, generator := range availableGenerators {
		if today.After(generator.availableSince) {
			selected = generator.factory
		}
	}

	return selected(manager, reportMonth)
}
