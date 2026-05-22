package blueticks

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// NewslettersResource exposes /v1/newsletters.
type NewslettersResource struct {
	client *Client
}

// Create creates a new WhatsApp newsletter (channel) via the connected engine.
//
// Create a new WhatsApp newsletter (channel). Requires `messages:write` scope
// (newsletter creation shares the messages write budget).
func (r *NewslettersResource) Create(ctx context.Context, params CreateNewsletterParams) (*Newsletter, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/newsletters",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	out := &Newsletter{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// List returns newsletters visible to the connected WhatsApp engine (cursor-paginated).
//
// List newsletters visible to the connected WhatsApp engine. Cursor-paginated
// via `limit` + `cursor`. Requires `newsletters:read` scope.
func (r *NewslettersResource) List(ctx context.Context, params ListNewslettersParams) (*NewslettersPage, error) {
	opts := requestOptions{method: http.MethodGet, path: "/v1/newsletters"}
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		q.Set("cursor", params.Cursor)
	}
	if len(q) > 0 {
		opts.query = q
	}
	body, err := r.client.do(ctx, opts)
	if err != nil {
		return nil, err
	}
	out := &NewslettersPage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Retrieve returns a newsletter by its JID.
//
// Retrieve a newsletter by its JID. Requires `newsletters:read` scope.
func (r *NewslettersResource) Retrieve(ctx context.Context, id string) (*Newsletter, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   fmt.Sprintf("/v1/newsletters/%s", id),
	})
	if err != nil {
		return nil, err
	}
	out := &Newsletter{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}
