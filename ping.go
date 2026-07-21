package blueticks

import (
	"context"
	"net/http"
)

// PingResource exposes /v1/ping.
type PingResource struct {
	client *Client
}

// Retrieve probes the API: returns liveness, the account id, and the currently
// connected WhatsApp engines for the authenticated API key. Useful as a
// connection test. No scope required.
func (r *PingResource) Retrieve(ctx context.Context) (*Ping, error) {
	body, err := r.client.do(ctx, requestOptions{method: http.MethodGet, path: "/v1/ping"})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data Ping `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}
