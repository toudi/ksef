package types

import monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

type Invoice struct {
	Attributes      map[string]string
	VAT             *VATInfo
	VATByRate       map[string]*VATInfo
	RegistryInvoice *monthlyregistry.Invoice
}
