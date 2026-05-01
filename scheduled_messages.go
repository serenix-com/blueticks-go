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

// List scheduled messages.
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

// Update a scheduled message.
//
// Updates the resource identified by id using data. Only the fields supplied
// on params are sent — unset (nil) fields preserve their server-side value.
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

// Delete (cancel) a scheduled message.
//
// Cancels a queued scheduled message before it fires. Returns the deleted
// ref ({id, deleted: true}).
func (r *ScheduledMessagesResource) Delete(ctx context.Context, id string) (*DeletedResource, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/v1/scheduled-messages/" + url.PathEscape(id),
	})
	if err != nil {
		return nil, err
	}
	out := &DeletedResource{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}
