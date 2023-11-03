package api

type BatchSession struct {
	api *API
}

func (a *API) BatchSessionInit() *BatchSession {
	return &BatchSession{api: a}
}
