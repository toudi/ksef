package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (a *APIClient) fullPath(path string) string {
	return fmt.Sprintf("https://%s%s", apiPrefix, path)
}

func (a *APIClient) NewRequest(ctx context.Context, method string, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, a.fullPath(path), nil)
}

func (a *APIClient) DoJSONResponse(req *http.Request, dest any) error {
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var decoder = json.NewDecoder(resp.Body)
	return decoder.Decode(dest)
}
