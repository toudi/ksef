package sessionregistry

func (r *Registry) PendingUploadHashes() []string {
	var hashes []string

	for _, session := range r.sessions {
		if session.Status == nil {
			for _, invoice := range session.Invoices {
				hashes = append(hashes, invoice.Checksum)
			}
		}
	}

	return hashes
}
