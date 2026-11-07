package monthly_test

import (
	"ksef/core/invoicedb/fs/metadata/registry/monthly"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvoiceID(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		type testCase struct {
			expected string
			input    monthly.Invoice
		}

		for _, test := range []testCase{
			{
				expected: "0001-f-123",
				input: monthly.Invoice{
					Type:   monthly.InvoiceTypeIssued,
					RefNo:  "f/123",
					OrdNum: 1,
				},
			},
			{
				expected: "0010-wystawca-f-456",
				input: monthly.Invoice{
					Type:   monthly.InvoiceTypeReceived,
					RefNo:  "f/456",
					OrdNum: 10,
					Issuer: &monthly.InvoiceIssuer{
						Name: "WYsTaWcA",
					},
				},
			},
			{
				expected: "0002-wystawca-f-789",
				input: monthly.Invoice{
					Type:   monthly.InvoiceTypePayer,
					RefNo:  "f/789",
					OrdNum: 2,
					Issuer: &monthly.InvoiceIssuer{
						Name: "WYsTaWcA",
					},
				},
			},
			{
				// Legacy invoices may have no issuer. ID must not panic.
				expected: "0003-f-999",
				input: monthly.Invoice{
					Type:   monthly.InvoiceTypeReceived,
					RefNo:  "f/999",
					OrdNum: 3,
				},
			},
		} {
			require.Equal(t, test.expected, test.input.ID())
		}
	})
}
