package registry

import (
	"ksef/internal/client/v2/types/invoices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistryAddInvoice(t *testing.T) {
	t.Run("make sure that we can replace invoice", func(t *testing.T) {
		const (
			origChecksum  = "87428fc522803d31065e7bce3cf03fe475096631e5e07bbd7a0fde60c4cf25c7" // sha256 of letter "a"
			newChecksum   = "0263829989b6fd954f72baaf2fc64bc2e2f01d692d4de72986ea808f6e99813f" // sha256 of letter "b"
			new2Checksum  = "a3a5e715f0cc574a73c3f9bebb6bc24f32ffd5b67b387244c2c909da779a1478" // sha256 of letter "c"
			invoiceRefNo  = "a/b/c"
			invoice2RefNo = "a/b/c/1"
		)
		var invoice = invoices.InvoiceMetadata{
			InvoiceNumber: invoiceRefNo,
			IssueDate:     "2025-12-31",
		}
		var invoice2 = invoices.InvoiceMetadata{
			InvoiceNumber: invoice2RefNo,
			IssueDate:     "2025-12-31",
		}
		r := NewRegistry()
		err := r.AddInvoice(
			invoice,
			origChecksum,
			nil,
		)
		require.NoError(t, err)
		require.Len(t, r.Invoices, 1)
		require.Equal(t, map[string]int{invoiceRefNo: 0}, r.refNoIndex)
		require.Equal(t, map[string]int{origChecksum: 0}, r.checksumIndex)
		require.Empty(t, r.seiRefNoIndex)
		// now let's add the same invoice but with a different checksum
		err = r.AddInvoice(
			invoice,
			newChecksum,
			nil,
		)
		require.NoError(t, err)
		// invoices array should not be extended
		require.Len(t, r.Invoices, 1)
		require.Equal(t, map[string]int{invoiceRefNo: 0}, r.refNoIndex)
		require.Equal(t, map[string]int{newChecksum: 0}, r.checksumIndex)
		require.Empty(t, r.seiRefNoIndex)
		// now let's add some different index
		err = r.AddInvoice(
			invoice2,
			new2Checksum,
			nil,
		)
		require.NoError(t, err)
		require.Len(t, r.Invoices, 2)
		require.Equal(t, map[string]int{invoiceRefNo: 0, invoice2RefNo: 1}, r.refNoIndex)
		require.Equal(t, map[string]int{newChecksum: 0, new2Checksum: 1}, r.checksumIndex)
		require.Empty(t, r.seiRefNoIndex)
	})

	t.Run("cannot replace invoice if existing one has KSeF number", func(t *testing.T) {
		const (
			origChecksum = "87428fc522803d31065e7bce3cf03fe475096631e5e07bbd7a0fde60c4cf25c7" // sha256 of letter "a"
			newChecksum  = "0263829989b6fd954f72baaf2fc64bc2e2f01d692d4de72986ea808f6e99813f" // sha256 of letter "b"
			invoiceRefNo = "a/b/c"
		)
		var invoice = invoices.InvoiceMetadata{
			InvoiceNumber: invoiceRefNo,
			IssueDate:     "2025-12-31",
			KSeFNumber:    "ksef-012",
		}
		var invoice2 = invoices.InvoiceMetadata{
			InvoiceNumber: invoiceRefNo,
			IssueDate:     "2025-12-31",
		}
		r := NewRegistry()
		err := r.AddInvoice(
			invoice,
			origChecksum,
			nil,
		)
		require.NoError(t, err)
		require.Len(t, r.Invoices, 1)
		require.Equal(t, map[string]int{invoiceRefNo: 0}, r.refNoIndex)
		require.Equal(t, map[string]int{origChecksum: 0}, r.checksumIndex)
		require.Equal(t, map[string]int{invoice.KSeFNumber: 0}, r.seiRefNoIndex)

		err = r.AddInvoice(
			invoice2,
			newChecksum,
			nil,
		)

		require.Equal(t, errConflictingRefsForSameChecksum, err)
		require.Len(t, r.Invoices, 1)
		require.Equal(t, map[string]int{invoiceRefNo: 0}, r.refNoIndex)
		require.Equal(t, map[string]int{origChecksum: 0}, r.checksumIndex)
		require.Equal(t, map[string]int{invoice.KSeFNumber: 0}, r.seiRefNoIndex)
	})
}
