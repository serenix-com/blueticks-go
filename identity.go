package blueticks

import (
	"context"
	"net/http"
	"net/url"
)

// IdentityResource exposes /v1/identity — resolving WhatsApp JIDs (@lid, @c.us,
// @s.whatsapp.net, @g.us) to canonical phone-form identities.
type IdentityResource struct {
	client *Client
}

// Get resolves a single WhatsApp JID to its canonical identity. An unknown
// identity returns Resolved=false (never a 404); a malformed id is a 400.
// Requires `contacts:read`.
func (r *IdentityResource) Get(ctx context.Context, id string) (*Identity, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/identity/" + url.PathEscape(id),
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data Identity `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// Resolve resolves up to 100 WhatsApp JIDs in one call. Order-preserving; each
// element matches the shape returned by Get. Requires `contacts:read`.
func (r *IdentityResource) Resolve(ctx context.Context, ids []string) (*ResolveIdentitiesResult, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/identity/resolve",
		body:   map[string][]string{"ids": ids},
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data []Identity `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &ResolveIdentitiesResult{Data: env.Data}, nil
}
