package api

import "net/url"

const EndpointAuthorisationChallenge = "online/Session/AuthorisationChallenge"
const EndpointInitToken = "online/Session/InitToken"
const EndpointLogout = "online/Session/Terminate"
const EndpointSendInvoice = "online/Invoice/Send"
const EndpointBatchInit = "batch/Init"
const EndpointBatchFinish = "batch/Finish"

func (a *API) apiEndpoint(path string) string {
	endpoint := &url.URL{Host: a.environment.host, Path: "/api/" + path, Scheme: "https"}
	return endpoint.String()
}
