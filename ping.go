package blueticks

import (
	"context"
	"net/http"
)

// PingResource exposes /v1/ping.
type PingResource struct {
	client *Client
}

// Retrieve probes the API: returns the account ID, key prefix, and granted scopes
// for the authenticated API key.
//
// Probe the API: returns the account ID, key prefix, and granted scopes for
// the authenticated API key. Useful as a connection test and to inspect what
// an integration is allowed to do. No scope required.
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
