package blueticks

import (
	"context"
	"net/http"
)

// PingResource exposes /v1/ping. The ping endpoint is a low-cost health
// check that also returns metadata about the calling API key (key prefix,
// scopes, account id) — useful for sanity-checking auth at boot time.
type PingResource struct {
	client *Client
}

// Retrieve returns key metadata for the authenticated API key.
func (r *PingResource) Retrieve(ctx context.Context) (*Ping, error) {
	body, err := r.client.do(ctx, requestOptions{method: http.MethodGet, path: "/v1/ping"})
	if err != nil {
		return nil, err
	}
	out := &Ping{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}
