package blueticks

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// WebhooksResource exposes /v1/webhooks.
type WebhooksResource struct {
	client *Client
}

// ListWebhooksParams is the query-param bundle for WebhooksResource.List.
type ListWebhooksParams struct {
	Limit  int
	Cursor string
}

// Create registers a new webhook subscription. The returned secret is
// shown only on creation — store it securely and use VerifyWebhook to
// validate incoming deliveries.
func (r *WebhooksResource) Create(ctx context.Context, params CreateWebhookParams) (*WebhookCreateResult, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/webhooks",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	out := &WebhookCreateResult{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// List webhooks, newest first. Cursor-paginated.
func (r *WebhooksResource) List(ctx context.Context, params ListWebhooksParams) (*WebhooksPage, error) {
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		q.Set("cursor", params.Cursor)
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/webhooks",
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	out := &WebhooksPage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get retrieves a webhook by id.
func (r *WebhooksResource) Get(ctx context.Context, id string) (*Webhook, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/webhooks/" + url.PathEscape(id),
	})
	if err != nil {
		return nil, err
	}
	out := &Webhook{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Update mutates one or more fields on a webhook subscription.
func (r *WebhooksResource) Update(ctx context.Context, id string, params UpdateWebhookParams) (*Webhook, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPatch,
		path:   "/v1/webhooks/" + url.PathEscape(id),
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	out := &Webhook{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete removes a webhook subscription.
func (r *WebhooksResource) Delete(ctx context.Context, id string) error {
	_, err := r.client.do(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/v1/webhooks/" + url.PathEscape(id),
	})
	return err
}

// RotateSecret invalidates the current signing secret and returns a new
// one. Use the new secret for all subsequent VerifyWebhook calls.
func (r *WebhooksResource) RotateSecret(ctx context.Context, id string) (*WebhookCreateResult, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/webhooks/" + url.PathEscape(id) + "/rotate-secret",
	})
	if err != nil {
		return nil, err
	}
	out := &WebhookCreateResult{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}
