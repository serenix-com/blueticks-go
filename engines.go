package blueticks

import (
	"context"
	"net/http"
)

// EnginesResource exposes /v1/engines.
type EnginesResource struct {
	client *Client
}

// Status reports the engine's current connection state.
func (r *EnginesResource) Status(ctx context.Context) (*EngineStatus, error) {
	body, err := r.client.do(ctx, requestOptions{method: http.MethodGet, path: "/v1/engines"})
	if err != nil {
		return nil, err
	}
	out := &EngineStatus{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}
