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

// Me returns the WhatsApp profile attached to the engine session.
func (r *EnginesResource) Me(ctx context.Context) (*WhatsAppMe, error) {
	body, err := r.client.do(ctx, requestOptions{method: http.MethodGet, path: "/v1/engines/me"})
	if err != nil {
		return nil, err
	}
	out := &WhatsAppMe{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Logout signs the engine out of WhatsApp Web (requires a fresh QR scan to
// reconnect).
func (r *EnginesResource) Logout(ctx context.Context) (map[string]any, error) {
	return r.action(ctx, "/v1/engines/logout")
}

// Reload restarts the engine session in-place. Use to clear transient
// stream wedges; preserves the current WhatsApp authentication.
func (r *EnginesResource) Reload(ctx context.Context) (map[string]any, error) {
	return r.action(ctx, "/v1/engines/reload")
}

func (r *EnginesResource) action(ctx context.Context, path string) (map[string]any, error) {
	body, err := r.client.do(ctx, requestOptions{method: http.MethodPost, path: path})
	if err != nil {
		return nil, err
	}
	out := map[string]any{}
	if err := decodeJSON(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}
