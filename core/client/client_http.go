package client

import "ksef/internal/http"

func (c *Client) authenticatedHTTPClient() *http.Client {
	return c.httpClient
}
