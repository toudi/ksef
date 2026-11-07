package monthly

import "ksef/core/invoicedb"

func (r *Registry) SyncParams() *invoicedb.InvoiceSyncParams {
	if r.state.SyncParams == nil {
		return nil
	}
	return &invoicedb.InvoiceSyncParams{
		LastTimestamp: r.state.SyncParams.LastTimestamp,
		SubjectTypes:  r.state.SyncParams.SubjectTypes,
	}
}

func (r *Registry) SetSyncParams(params invoicedb.InvoiceSyncParams) {
	if r.state.SyncParams == nil {
		r.state.SyncParams = &SyncParams{}
	}
	r.state.SyncParams.LastTimestamp = params.LastTimestamp
	if params.SubjectTypes != nil {
		r.state.SyncParams.SubjectTypes = params.SubjectTypes
	}
	r.dirty = true
}
