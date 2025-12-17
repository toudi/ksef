package sessionregistry

import (
	sessionTypes "ksef/internal/client/v2/session/types"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/utils"
)

func (r *Registry) Update(
	uploadResult []*sessionTypes.UploadSessionResult,
	invoiceChecksumToRegistry map[string]*monthlyregistry.Registry,
) error {
	return r.processUploadResult(
		uploadResult,
		invoiceChecksumToRegistry,
	)
}

func (r *Registry) processUploadResult(
	uploadResult []*sessionTypes.UploadSessionResult,
	invoiceChecksumToRegistry map[string]*monthlyregistry.Registry,
) error {
	var affectedRegistries = make(map[*monthlyregistry.Registry]bool)

	// first step - update entries
	for _, uploadSessionStatus := range uploadResult {
		uploadSessionId := uploadSessionStatus.SessionID
		entry, entryIndex, exists := r.lookupSessionById(uploadSessionId)
		if !exists {
			entry = &UploadSession{
				Timestamp: uploadSessionStatus.Timestamp,
				RefNo:     uploadSessionId,
				Invoices:  make([]*Invoice, 0, len(uploadSessionStatus.Invoices)),
			}
		}

		if uploadSessionStatus.Status != nil && entry.Status == nil {
			entry.Status = uploadSessionStatus.Status
		}

		for _, invoiceUploadStatus := range uploadSessionStatus.Invoices {
			r.logger.Debug("update monthly registry with", "upload session status", invoiceUploadStatus)

			if invoiceUploadStatus.Status.Successful() {
				var invoiceChecksum = utils.Base64ToHex(invoiceUploadStatus.Checksum)
				// unfortunetely, KSeF API returns checksums in base64 form which is space-efficient,
				// but makes it painful to compare with local checksums generated with command-line
				// utilities, therefore internally we're keeping the hex-encoded checksum
				var registry = invoiceChecksumToRegistry[invoiceChecksum]
				// represents invoice original ref no, extracted from the registry
				if err := registry.UpdateInvoiceByChecksum(
					invoiceChecksum,
					func(invoice *monthlyregistry.Invoice) {
						// basically, the reference number is assigned during the interactive session even
						// before the invoice is processed and therefore it's essential to know
						// if it was processed succesfully to update the reference number in the registry.
						// otherwise, we'd perpetually replace the invoice's KSeF number
						if invoiceUploadStatus.KSeFRefNo != "" && invoiceUploadStatus.Status.Successful() {
							invoice.KSeFRefNo = invoiceUploadStatus.KSeFRefNo
						} else {
							if len(invoiceUploadStatus.Status.Details) > 0 {
								invoice.UploadErrors = invoiceUploadStatus.Status.Details
							}
						}
					},
				); err != nil {
					return err
				}

				// if we're here then the updateInvoiceByChecksum succeeded and
				// we've now successfully obtained invoice's original number as
				// an artifact
				affectedRegistries[registry] = true
			}

			entry.addInfoAboutInvoice(
				invoiceUploadStatus,
			)
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
