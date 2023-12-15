package client

import (
	"fmt"
	"net/url"
)

func (a *APIClient) apiEndpoint(path string) string {
	endpoint, _ := url.Parse(fmt.Sprintf("https://%s/api/"+path, a.Environment.Host))
	// endpoint := &url.URL{Host: a.Environment.Host, Path: "/api/" + path, Scheme: "https"}
	return endpoint.String()
}
