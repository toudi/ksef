package registry

import (
	"ksef/internal/sei/api/client/v2/types/invoices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOutputFilename(t *testing.T) {
	type testCase struct {
		subjectType invoices.SubjectType
		meta        invoices.InvoiceMetadata
		expected    string
	}

	var registry = &InvoiceRegistry{
		sourcePath: "/tmp/test.yaml",
	}

	for _, test := range []testCase{
		{
			subjectType: invoices.SubjectTypeIssuer,
			meta: invoices.InvoiceMetadata{
				Seller: invoices.InvoiceSubjectMetadata{
					NIP: "sellerId123",
				},
				InvoiceNumber: "a/b/c",
			},
			expected: "/tmp/wystawione/001-sellerid123-a-b-c.xml",
		},
	} {
		require.Equal(t, test.expected, registry.GetTargetFilename(test.meta, test.subjectType))
	}
}
