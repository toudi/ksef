package sessionregistry

func (r *Registry) lookupSessionById(uploadSessionId string) (*UploadSession, int, bool) {
	for index, entry := range r.sessions {
		if entry.RefNo == uploadSessionId {
			return entry, index, true
		}
	}

	return nil, -1, false
}
