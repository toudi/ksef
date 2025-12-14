package sessionregistry

import (
	sessionTypes "ksef/internal/client/v2/session/types"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
)

func (r *Registry) Update(
	uploadResult []*sessionTypes.UploadSessionResult,
	invoiceChecksumToRegistry map[string]*monthlyregistry.Registry,
) error {
	return r.processUploadResult(
		uploadResult,
		invoiceChecksumToRegistry,
		func(registry *monthlyregistry.Registry, invoiceStatus sessionTypes.InvoiceUploadResult) error {
			return registry.UpdateInvoiceByChecksum(
				invoiceStatus.Checksum,
				func(invoice *monthlyregistry.Invoice) {
					if invoiceStatus.KSeFRefNo != "" {
						invoice.KSeFRefNo = invoiceStatus.KSeFRefNo
					}
					invoice.UploadErrors = invoiceStatus.Errors
				},
			)
		},
	)
}

func (r *Registry) UpdateUploadedInvoicesResult(
	uploadResult []*sessionTypes.UploadSessionResult,
	invoiceChecksumToRegistry map[string]*monthlyregistry.Registry,
) error {
	return r.processUploadResult(
		uploadResult,
		invoiceChecksumToRegistry,
		func(registry *monthlyregistry.Registry, invoiceStatus sessionTypes.InvoiceUploadResult) error {
			return registry.UpdateInvoiceByChecksum(
				invoiceStatus.Checksum,
				func(invoice *monthlyregistry.Invoice) {
					if len(invoiceStatus.Errors) > 0 {
						invoice.UploadErrors = invoiceStatus.Errors
					}
				},
			)
		},
	)
}

func (r *Registry) processUploadResult(
	uploadResult []*sessionTypes.UploadSessionResult,
	invoiceChecksumToRegistry map[string]*monthlyregistry.Registry,
	updateHandler func(
		registry *monthlyregistry.Registry,
		invoiceStatus sessionTypes.InvoiceUploadResult,
	) error,
) error {
	var affectedRegistries = make(map[*monthlyregistry.Registry]bool)

	// first step - update entries
	for _, uploadSessionStatus := range uploadResult {
		uploadSessionId := uploadSessionStatus.SessionID
		entry, entryIndex, exists := r.lookupSessionById(uploadSessionId)
		if !exists {
			entry = &UploadSession{
				RefNo:    uploadSessionId,
				Invoices: make([]*Invoice, 0, len(uploadSessionStatus.Invoices)),
			}
		}

		for _, invoiceUploadStatus := range uploadSessionStatus.Invoices {
			var registry = invoiceChecksumToRegistry[invoiceUploadStatus.Checksum]
			if err := updateHandler(registry, invoiceUploadStatus); err != nil {
				return err
			}
			affectedRegistries[registry] = true
		}

		if uploadSessionStatus.Status != nil && len(uploadSessionStatus.Status.Upo.Pages) > 0 {
			entry.UPO = uploadSessionStatus.Status.Upo.Pages
		}

		if exists {
			r.sessions[entryIndex] = entry
		} else {
			r.sessions = append(r.sessions, entry)
		}
		r.dirty = true
	}
	// second step - persist the file.
	if err := r.Save(); err != nil {
		return err
	}

	// third step - iterate through affected registries and save their files
	for registry := range affectedRegistries {
		if err := registry.Save(); err != nil {
			return err
		}
	}

	return nil
}
