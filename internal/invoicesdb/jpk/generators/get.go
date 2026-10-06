package generators

import (
	"time"

	"ksef/internal/invoicesdb/annotations"
	"ksef/internal/invoicesdb/jpk/generators/interfaces"
	"ksef/internal/invoicesdb/jpk/generators/jpk_v7m_3"
	subjectsettings "ksef/internal/invoicesdb/subject-settings"
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
	annotations *annotations.Annotations,
	subjectSettings *subjectsettings.SubjectSettings,
	reportMonth time.Time,
) interfaces.JPKGenerator {
	var selected interfaces.JPKGeneratorFactory
	today := time.Now().Local()

	for _, generator := range availableGenerators {
		if today.After(generator.availableSince) {
			selected = generator.factory
		}
	}

	return selected(annotations, subjectSettings, reportMonth)
}
