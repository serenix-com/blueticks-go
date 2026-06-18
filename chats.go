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

// ListParticipants returns chat participants. For group chats, the
// participant list is paginated; for DMs, the single counterparty is
// returned in a one-element page.
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

// MarkRead clears the unread badge on the connected engine for the given chat.
func (r *ChatsResource) MarkRead(ctx context.Context, chatID string) (*OkResponse, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/chats/" + url.PathEscape(chatID) + "/mark_read",
	})
	if err != nil {
		return nil, err
	}
	out := &OkResponse{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Open brings the chat to the foreground on the connected WhatsApp Web client
// (creating the chat if it doesn't exist yet for the engine). Useful before
// issuing follow-up reads on a fresh JID.
func (r *ChatsResource) Open(ctx context.Context, chatID string) (*ChatRef, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/chats/" + url.PathEscape(chatID) + "/open",
	})
	if err != nil {
		return nil, err
	}
	out := &ChatRef{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// SendMessage sends a message immediately to a specific chat.
//
// The body is the same discriminated union as POST /v1/scheduled-messages minus
// `to` (derived from the URL path) and `sendAt` (this endpoint is
// fire-and-forget). Variants: "text", "media", "poll". The dispatch is direct
// — no DB row is created; the response carries the WhatsApp wire key under
// `key`. For scheduled or queue-managed sends use ScheduledMessagesResource.Create
// instead. Pass an IdempotencyKey on the params to make the call safe to retry.
func (r *ChatsResource) SendMessage(ctx context.Context, chatID string, params SendMessageInChatParams) (*ScheduledMessage, error) {
	body, err := r.client.do(ctx, requestOptions{
		method:         http.MethodPost,
		path:           "/v1/chats/" + url.PathEscape(chatID) + "/messages",
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

// ListMessages returns messages from a chat.
//
// `MessageTypes` is sent as a single comma-separated query value because the
// OpenAPI schema declares `style: form, explode: false` for that param.
func (r *ChatsResource) ListMessages(ctx context.Context, chatID string, params ListChatMessagesParams) (*ChatMessagesPage, error) {
	q := url.Values{}
	if params.Order != "" {
		q.Set("order", params.Order)
	}
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
		q.Set("messageTypes", strings.Join(params.MessageTypes, ","))
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

// GetMessage retrieves a single message by its complete WhatsApp message key.
//
// Pass an empty chatID to omit the optional chatId query param (the engine
// resolves the chat from the key when possible).
func (r *ChatsResource) GetMessage(ctx context.Context, waMessageKey, chatID string) (*ChatMessage, error) {
	q := url.Values{}
	if chatID != "" {
		q.Set("chatId", chatID)
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/messages/" + url.PathEscape(waMessageKey),
		query:  q,
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

// GetMessageAck returns the WhatsApp ack value for a sent message:
// -1=error, 0=pending, 1=server, 2=device, 3=read, 4=played; nil when no
// engine response.
//
// Pass an empty chatID to omit the optional chatId query param.
func (r *ChatsResource) GetMessageAck(ctx context.Context, waMessageKey, chatID string) (*MessageAck, error) {
	q := url.Values{}
	if chatID != "" {
		q.Set("chatId", chatID)
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/messages/ack/" + url.PathEscape(waMessageKey),
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	out := &MessageAck{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// React adds or replaces a reaction emoji on a message. Pass an empty
// emoji to clear an existing reaction.
//
// Pass an empty chatID to omit the optional chatId query param.
func (r *ChatsResource) React(ctx context.Context, waMessageKey, chatID, emoji string) (*OkResponse, error) {
	q := url.Values{}
	if chatID != "" {
		q.Set("chatId", chatID)
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/messages/reactions/" + url.PathEscape(waMessageKey),
		query:  q,
		body:   map[string]string{"emoji": emoji},
	})
	if err != nil {
		return nil, err
	}
	out := &OkResponse{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Pin pins a message to the top of its chat by its complete WhatsApp message key.
// duration is the pin expiry in seconds; pass 0 to use WhatsApp's default of 7 days.
// Pass an empty chatID to omit the optional chatId query param.
func (r *ChatsResource) Pin(ctx context.Context, waMessageKey, chatID string, duration int) (*OkResponse, error) {
	q := url.Values{}
	if chatID != "" {
		q.Set("chatId", chatID)
	}
	var reqBody interface{}
	if duration > 0 {
		reqBody = map[string]int{"duration": duration}
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/messages/pin/" + url.PathEscape(waMessageKey),
		query:  q,
		body:   reqBody,
	})
	if err != nil {
		return nil, err
	}
	out := &OkResponse{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Unpin removes an existing pin from a message by its complete WhatsApp message key.
// Pass an empty chatID to omit the optional chatId query param.
func (r *ChatsResource) Unpin(ctx context.Context, waMessageKey, chatID string) (*OkResponse, error) {
	q := url.Values{}
	if chatID != "" {
		q.Set("chatId", chatID)
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/messages/unpin/" + url.PathEscape(waMessageKey),
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	out := &OkResponse{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// LoadOlderMessages asks the engine to pull older history from the connected
// phone for chats that haven't been fully synced yet. Use this once before
// paginating with `since` if you need messages older than what's already
// cached.
func (r *ChatsResource) LoadOlderMessages(ctx context.Context, chatID string) (*LoadOlderMessagesResponse, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/chats/" + url.PathEscape(chatID) + "/messages/load_older",
	})
	if err != nil {
		return nil, err
	}
	out := &LoadOlderMessagesResponse{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMedia downloads message media. May return base64 inline depending on
// configuration / media kind; check ChatMedia.MediaUnavailable for failure
// kinds.
//
// Pass an empty chatID to omit the optional chatId query param.
func (r *ChatsResource) GetMedia(ctx context.Context, waMessageKey, chatID string) (*ChatMedia, error) {
	q := url.Values{}
	if chatID != "" {
		q.Set("chatId", chatID)
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/messages/media/" + url.PathEscape(waMessageKey),
		query:  q,
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

// GetMediaURL returns a hosted URL for the message media without inlining
// bytes. Faster + cheaper than GetMedia when the caller can fetch the URL
// themselves.
//
// Pass an empty chatID to omit the optional chatId query param.
func (r *ChatsResource) GetMediaURL(ctx context.Context, waMessageKey, chatID string) (*MediaURLResponse, error) {
	q := url.Values{}
	if chatID != "" {
		q.Set("chatId", chatID)
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/messages/media_url/" + url.PathEscape(waMessageKey),
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	out := &MediaURLResponse{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// BatchMessageAcks fetches ack state for up to 200 sent messages in one call.
// Useful for campaign dashboards and status reconciliation.
func (r *ChatsResource) BatchMessageAcks(ctx context.Context, params BatchMessageAcksParams) (*BatchMessageAcksResponse, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/chats/message_acks",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	out := &BatchMessageAcksResponse{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}
