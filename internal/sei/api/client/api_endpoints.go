package client

import "net/url"

func (a *APIClient) apiEndpoint(path string) string {
	endpoint := &url.URL{Host: a.Environment.Host, Path: "/api/" + path, Scheme: "https"}
	return endpoint.String()
}
