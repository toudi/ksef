package registry

import "ksef/internal/logging"

func (r *InvoiceRegistry) AddUPOToSession(sessionId string, upoRefNo string) {
	if _, exists := r.UploadSessions[sessionId]; !exists {
		logging.UPOLogger.With("session id", sessionId, "upo", upoRefNo).Error("próba przypisania UPO do nieistniejącej sesji")
		return
	}
	r.UploadSessions[sessionId].UPO = append(r.UploadSessions[sessionId].UPO, upoRefNo)
}
