package blueticks

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// MessagesResource exposes /v1/messages.
type MessagesResource struct {
	client *Client
}

// ListMessagesParams is the query-param bundle for MessagesResource.List.
// Both fields are optional; pass the zero value to omit them.
type ListMessagesParams struct {
	Limit  int
	Cursor string
}

// Send (or schedule) a message.
//
// Pass an IdempotencyKey on the request to make this call safe to retry. The
// key is sent as the `Idempotency-Key` header — the server short-circuits
// duplicate calls within its replay window.
func (r *MessagesResource) Send(ctx context.Context, params SendMessageParams) (*Message, error) {
	body, err := r.client.do(ctx, requestOptions{
		method:         http.MethodPost,
		path:           "/v1/messages",
		body:           params,
		idempotencyKey: params.IdempotencyKey,
	})
	if err != nil {
		return nil, err
	}
	out := &Message{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get retrieves a message by id.
func (r *MessagesResource) Get(ctx context.Context, id string) (*Message, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/messages/" + url.PathEscape(id),
	})
	if err != nil {
		return nil, err
	}
	out := &Message{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// List messages sent through the API, newest first. Cursor-paginated.
func (r *MessagesResource) List(ctx context.Context, params ListMessagesParams) (*MessagesPage, error) {
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		q.Set("cursor", params.Cursor)
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/messages",
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	out := &MessagesPage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}
