package blueticks

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// ScheduledMessagesResource exposes /v1/scheduled-messages.
type ScheduledMessagesResource struct {
	client *Client
}

// Send message.
//
// Creates (and optionally schedules) a message. Pass an IdempotencyKey on
// the request to make this call safe to retry — the key is sent as the
// `Idempotency-Key` header and the server short-circuits duplicate calls
// within its replay window.
func (r *ScheduledMessagesResource) Create(ctx context.Context, params CreateScheduledMessageParams) (*ScheduledMessage, error) {
	body, err := r.client.do(ctx, requestOptions{
		method:         http.MethodPost,
		path:           "/v1/scheduled-messages",
		body:           params,
		idempotencyKey: params.IdempotencyKey,
	})
	if err != nil {
		return nil, err
	}
	out := &ScheduledMessage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// List messages.
//
// Retrieves a list of all resources from the service. Cursor-paginated.
func (r *ScheduledMessagesResource) List(ctx context.Context, params ListScheduledMessagesParams) (*ScheduledMessagesPage, error) {
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		q.Set("cursor", params.Cursor)
	}
	if params.ChatID != "" {
		q.Set("chatId", params.ChatID)
	}
	if params.Status != "" {
		q.Set("status", params.Status)
	}
	if params.Q != "" {
		q.Set("q", params.Q)
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/scheduled-messages",
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	out := &ScheduledMessagesPage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get message.
//
// Retrieves a single resource with the given id from the service.
func (r *ScheduledMessagesResource) Retrieve(ctx context.Context, id string) (*ScheduledMessage, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/scheduled-messages/" + url.PathEscape(id),
	})
	if err != nil {
		return nil, err
	}
	out := &ScheduledMessage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Update message.
//
// Updates the resource identified by id using the supplied fields. Only the
// fields supplied on params are sent — unset (nil) fields preserve their
// server-side value. Rejected once the message has advanced past the
// editable window (status not in pending/sending).
func (r *ScheduledMessagesResource) Update(ctx context.Context, id string, params UpdateScheduledMessageParams) (*ScheduledMessage, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPatch,
		path:   "/v1/scheduled-messages/" + url.PathEscape(id),
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	out := &ScheduledMessage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}
