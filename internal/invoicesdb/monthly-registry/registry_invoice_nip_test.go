package monthlyregistry

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistryGetNip(t *testing.T) {
	reg := &Registry{
		dir: "/foo/bar/ksef-test.mf.gov.pl/1112223344/2026/01",
	}

	require.Equal(t, "1112223344", reg.GetNIP())
}
