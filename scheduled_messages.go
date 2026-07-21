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

// Create (schedule) a message to a chat.
//
// POSTs to /v1/scheduled-messages/{chatId}. Set params.SendAt to schedule for
// later; omit for immediate queueing. Pass an IdempotencyKey on params to make
// the call safe to retry — it is sent as the Idempotency-Key header.
func (r *ScheduledMessagesResource) Create(ctx context.Context, chatID string, params SendMessageParams) (*SentMessage, error) {
	body, err := r.client.do(ctx, requestOptions{
		method:         http.MethodPost,
		path:           "/v1/scheduled-messages/" + url.PathEscape(chatID),
		body:           params,
		idempotencyKey: params.IdempotencyKey,
	})
	if err != nil {
		return nil, err
	}
	return decodeSentMessage(body)
}

// List scheduled messages. Offset-paginated.
func (r *ScheduledMessagesResource) List(ctx context.Context, params ListScheduledMessagesParams) (*ScheduledMessagesPage, error) {
	q := url.Values{}
	if params.ChatID != "" {
		q.Set("chatId", params.ChatID)
	}
	if params.SearchToken != "" {
		q.Set("searchToken", params.SearchToken)
	}
	if params.Status != "" {
		q.Set("status", params.Status)
	}
	if params.Order != "" {
		q.Set("order", params.Order)
	}
	if params.Skip > 0 {
		q.Set("skip", strconv.Itoa(params.Skip))
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
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

// Retrieve a single scheduled message by id.
func (r *ScheduledMessagesResource) Retrieve(ctx context.Context, id string) (*SentMessage, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/scheduled-messages/" + url.PathEscape(id),
	})
	if err != nil {
		return nil, err
	}
	return decodeSentMessage(body)
}

// Update mutates a queued/scheduled message. Only supplied (non-nil) fields are
// sent. Rejected once the message has advanced past the editable window.
func (r *ScheduledMessagesResource) Update(ctx context.Context, id string, params UpdateScheduledMessageParams) (*SentMessage, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPatch,
		path:   "/v1/scheduled-messages/" + url.PathEscape(id),
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	return decodeSentMessage(body)
}

// Delete cancels a scheduled message. Returns {id, deleted: true}.
func (r *ScheduledMessagesResource) Delete(ctx context.Context, id string) (*DeletedResource, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/v1/scheduled-messages/" + url.PathEscape(id),
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data DeletedResource `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

func decodeSentMessage(body []byte) (*SentMessage, error) {
	var env struct {
		Data SentMessage `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}
