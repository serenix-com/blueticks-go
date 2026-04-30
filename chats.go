package blueticks

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ChatsResource exposes /v1/chats and /v1/chats/{id}/* — every method is
// a thin wrapper that dispatches to the user's WhatsApp engine.
type ChatsResource struct {
	client *Client
}

// ListChatsParams is the query-param bundle for ChatsResource.List.
type ListChatsParams struct {
	Query  string
	Limit  int
	Cursor string
}

// ListParticipantsParams is the query-param bundle for ChatsResource.ListParticipants.
type ListParticipantsParams struct {
	Limit  int
	Cursor string
}

// List/search chats, newest first. Cursor-paginated.
func (r *ChatsResource) List(ctx context.Context, params ListChatsParams) (*ChatsPage, error) {
	q := url.Values{}
	if params.Query != "" {
		q.Set("query", params.Query)
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		q.Set("cursor", params.Cursor)
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/chats",
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	out := &ChatsPage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get retrieves a chat by its JID.
func (r *ChatsResource) Get(ctx context.Context, chatID string) (*Chat, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/chats/" + url.PathEscape(chatID),
	})
	if err != nil {
		return nil, err
	}
	out := &Chat{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListParticipants in a group chat. Cursor-paginated.
func (r *ChatsResource) ListParticipants(ctx context.Context, chatID string, params ListParticipantsParams) (*ParticipantsPage, error) {
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		q.Set("cursor", params.Cursor)
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/chats/" + url.PathEscape(chatID) + "/participants",
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	out := &ParticipantsPage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// MarkRead marks a chat as read.
func (r *ChatsResource) MarkRead(ctx context.Context, chatID string) (map[string]any, error) {
	return r.actionMap(ctx, chatID, "mark_read")
}

// Open the chat in the engine (signals presence/read receipts).
func (r *ChatsResource) Open(ctx context.Context, chatID string) (map[string]any, error) {
	return r.actionMap(ctx, chatID, "open")
}

func (r *ChatsResource) actionMap(ctx context.Context, chatID, action string) (map[string]any, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/chats/" + url.PathEscape(chatID) + "/" + action,
	})
	if err != nil {
		return nil, err
	}
	out := map[string]any{}
	if err := decodeJSON(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListMessages in a chat.
//
// `MessageTypes` is sent as a single comma-separated query value because the
// OpenAPI schema declares `style: form, explode: false` for that param.
func (r *ChatsResource) ListMessages(ctx context.Context, chatID string, params ListChatMessagesParams) (*ChatMessagesPage, error) {
	q := url.Values{}
	mode := params.Mode
	if mode == "" {
		mode = "latest"
	}
	q.Set("mode", mode)
	if params.Query != "" {
		q.Set("query", params.Query)
	}
	if params.Since != "" {
		q.Set("since", params.Since)
	}
	if params.Until != "" {
		q.Set("until", params.Until)
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		q.Set("cursor", params.Cursor)
	}
	if len(params.MessageTypes) > 0 {
		q.Set("message_types", strings.Join(params.MessageTypes, ","))
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/chats/" + url.PathEscape(chatID) + "/messages",
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	out := &ChatMessagesPage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMessage retrieves a single message by WhatsApp message key.
func (r *ChatsResource) GetMessage(ctx context.Context, chatID, key string) (*ChatMessage, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/chats/" + url.PathEscape(chatID) + "/messages/" + url.PathEscape(key),
	})
	if err != nil {
		return nil, err
	}
	out := &ChatMessage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMessageAck returns the per-recipient acknowledgement state for a
// single message. Shape varies by chat type, so the response is returned
// as a free-form map for forward compatibility.
func (r *ChatsResource) GetMessageAck(ctx context.Context, chatID, key string) (MessageAcksResponse, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/chats/" + url.PathEscape(chatID) + "/messages/" + url.PathEscape(key) + "/ack",
	})
	if err != nil {
		return nil, err
	}
	out := MessageAcksResponse{}
	if err := decodeJSON(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// React adds or replaces a reaction emoji on a message. Pass an empty
// emoji to clear an existing reaction.
func (r *ChatsResource) React(ctx context.Context, chatID, key, emoji string) (map[string]any, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/chats/" + url.PathEscape(chatID) + "/messages/" + url.PathEscape(key) + "/reactions",
		body:   map[string]string{"emoji": emoji},
	})
	if err != nil {
		return nil, err
	}
	out := map[string]any{}
	if err := decodeJSON(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// LoadOlderMessages tells the engine to fetch and persist older history
// for a chat. Returns when the load is complete (or the engine reports
// no further history is available).
func (r *ChatsResource) LoadOlderMessages(ctx context.Context, chatID string) (map[string]any, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/chats/" + url.PathEscape(chatID) + "/messages/load_older",
	})
	if err != nil {
		return nil, err
	}
	out := map[string]any{}
	if err := decodeJSON(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMedia downloads message media. May return base64 inline depending on
// configuration / media kind; check ChatMedia.MediaUnavailable for failure
// kinds.
func (r *ChatsResource) GetMedia(ctx context.Context, chatID, key string) (*ChatMedia, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/chats/" + url.PathEscape(chatID) + "/messages/" + url.PathEscape(key) + "/media",
	})
	if err != nil {
		return nil, err
	}
	out := &ChatMedia{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMediaURL returns a short-lived signed URL pointing at the media for a message.
func (r *ChatsResource) GetMediaURL(ctx context.Context, chatID, key string) (map[string]any, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/chats/" + url.PathEscape(chatID) + "/messages/" + url.PathEscape(key) + "/media_url",
	})
	if err != nil {
		return nil, err
	}
	out := map[string]any{}
	if err := decodeJSON(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BatchMessageAcks fetches ack state for many messages in one call.
// The response is keyed by message key.
func (r *ChatsResource) BatchMessageAcks(ctx context.Context, messageKeys []string) (MessageAcksResponse, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/chats/message_acks",
		body:   map[string]any{"message_keys": messageKeys},
	})
	if err != nil {
		return nil, err
	}
	out := MessageAcksResponse{}
	if err := decodeJSON(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}
