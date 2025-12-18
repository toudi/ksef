package sessionregistry

func (r *Registry) PendingUploadSessions() (sessions []*UploadSession) {
	for _, session := range r.sessions {
		if session.Status == nil {
			sessions = append(sessions, session)
		}
	}

	return sessions
}
