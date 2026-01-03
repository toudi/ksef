package jpk

import (
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	subjectsettings "ksef/internal/invoicesdb/subject-settings"
	"path/filepath"

	"github.com/spf13/viper"
)

type JPKManager struct {
	vip      *viper.Viper
	registry *monthlyregistry.Registry
	ss       *subjectsettings.SubjectSettings
}

func Manager(vip *viper.Viper, initializers ...func(*JPKManager) error) (*JPKManager, error) {
	manager := &JPKManager{}

	for _, initializer := range initializers {
		if err := initializer(manager); err != nil {
			return nil, err
		}
	}

	return manager, nil
}

func WithMonthlyRegistry(registry *monthlyregistry.Registry) func(*JPKManager) error {
	return func(m *JPKManager) error {
		m.registry = registry

		// we can load subject settings (if they exist) based on data from monthly
		// registry
		var err error
		m.ss, err = subjectsettings.OpenOrCreate(
			filepath.Join(
				registry.Dir(), // e.g. 2025/12
				"..",           // e.g. 2025
				"..",           // nip-level
			),
		)
		if err != nil {
			// log error
		}

		return nil
	}
}
