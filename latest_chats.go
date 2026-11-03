package blueticks

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// LatestChatsResource exposes /v1/latest-chats.
type LatestChatsResource struct {
	client *Client
}

// List returns the most-recent non-archived chats (pinned first), each with its
// latest messages (newest first), in one call — no per-chat round-trips.
// Requires `chats:read`.
func (r *LatestChatsResource) List(ctx context.Context, params LatestChatsParams) (*LatestChatsResult, error) {
	q := url.Values{}
	if params.NumberOfChats != nil {
		q.Set("numberOfChats", strconv.Itoa(*params.NumberOfChats))
	}
	if params.NumberOfMessages != nil {
		q.Set("numberOfMessages", strconv.Itoa(*params.NumberOfMessages))
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/latest-chats",
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data LatestChatsResult `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}
