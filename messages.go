package blueticks

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// MessagesResource exposes /v1/messages — reading chat history, sending
// messages, delivery acks, media, pins, and reactions. Every method dispatches
// to the connected WhatsApp engine.
type MessagesResource struct {
	client *Client
}

// List messages across all chats (or scoped to ChatID). Offset-paginated.
// MessageTypes is sent as a single comma-separated query value.
func (r *MessagesResource) List(ctx context.Context, params ListMessagesParams) (*MessagesPage, error) {
	q := url.Values{}
	if params.ChatID != "" {
		q.Set("chatId", params.ChatID)
	}
	if params.SearchToken != "" {
		q.Set("searchToken", params.SearchToken)
	}
	if params.Order != "" {
		q.Set("order", params.Order)
	}
	if params.Since != "" {
		q.Set("since", params.Since)
	}
	if params.Until != "" {
		q.Set("until", params.Until)
	}
	if len(params.MessageTypes) > 0 {
		q.Set("messageTypes", strings.Join(params.MessageTypes, ","))
	}
	for _, term := range params.QueryAny {
		q.Add("queryAny", term)
	}
	if params.Mime != "" {
		q.Set("mime", params.Mime)
	}
	if params.FilenameContains != "" {
		q.Set("filenameContains", params.FilenameContains)
	}
	if params.HasMedia != nil {
		q.Set("hasMedia", strconv.FormatBool(*params.HasMedia))
	}
	if params.Sender != "" {
		q.Set("sender", params.Sender)
	}
	if params.LoadFromPhoneIfNeeded != nil {
		q.Set("loadFromPhoneIfNeeded", strconv.FormatBool(*params.LoadFromPhoneIfNeeded))
	}
	if params.IncludeMediaContent != nil {
		q.Set("includeMediaContent", strconv.FormatBool(*params.IncludeMediaContent))
	}
	if params.Skip > 0 {
		q.Set("skip", strconv.Itoa(params.Skip))
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
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

// Get retrieves a single message by its WhatsApp message key. params.ChatID is
// an optional locator hint.
func (r *MessagesResource) Get(ctx context.Context, waMessageKey string, params GetMessageParams) (*Message, error) {
	q := url.Values{}
	if params.ChatID != "" {
		q.Set("chatId", params.ChatID)
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/messages/" + url.PathEscape(waMessageKey),
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data Message `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// Send a message immediately to a chat (POST /v1/messages/{chatId}).
//
// The body is a discriminated union on Type ("text", "media", "poll"). The
// dispatch is direct — no queue row is created. Pass an IdempotencyKey on
// params to make the call safe to retry.
func (r *MessagesResource) Send(ctx context.Context, chatID string, params SendInChatParams) (*SentMessage, error) {
	body, err := r.client.do(ctx, requestOptions{
		method:         http.MethodPost,
		path:           "/v1/messages/" + url.PathEscape(chatID),
		body:           params,
		idempotencyKey: params.IdempotencyKey,
	})
	if err != nil {
		return nil, err
	}
	return decodeSentMessage(body)
}

// GetAck returns the WhatsApp ack for a sent message: -1=error, 0=pending,
// 1=server, 2=device, 3=read, 4=played; nil when no engine response.
func (r *MessagesResource) GetAck(ctx context.Context, waMessageKey string, params MessageAckParams) (*MessageAck, error) {
	q := url.Values{}
	if params.ChatID != "" {
		q.Set("chatId", params.ChatID)
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/messages/ack/" + url.PathEscape(waMessageKey),
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data MessageAck `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// BatchAcks fetches ack state for up to 200 message keys in one call.
func (r *MessagesResource) BatchAcks(ctx context.Context, params BatchMessageAcksParams) (*BatchMessageAcksResult, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/messages/acks",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	out := &BatchMessageAcksResult{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// LoadOlder asks the engine to pull older history from the connected phone for a
// chat that isn't fully synced yet.
//
// To pull several pages in one call or stop at a target depth, use
// LoadOlderWithParams.
func (r *MessagesResource) LoadOlder(ctx context.Context, chatID string) (*LoadOlderResult, error) {
	return r.LoadOlderWithParams(ctx, chatID, LoadOlderParams{})
}

// LoadOlderWithParams is LoadOlder with the optional Pages (1-10, WA-paced
// ~2.5s apart) and UntilDate controls. The result reports how many pages ran
// (PagesExecuted) and how far back the chat is now loaded (Coverage).
func (r *MessagesResource) LoadOlderWithParams(ctx context.Context, chatID string, params LoadOlderParams) (*LoadOlderResult, error) {
	var reqBody interface{}
	if params.Pages != nil || params.UntilDate != nil {
		reqBody = params
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/messages/load_older/" + url.PathEscape(chatID),
		body:   reqBody,
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data LoadOlderResult `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// GetMedia downloads a message's media (inline base64 or hosted URL). Check
// Media.MediaUnavailable for failure kinds.
func (r *MessagesResource) GetMedia(ctx context.Context, waMessageKey string, params GetMediaParams) (*Media, error) {
	q := url.Values{}
	if params.ChatID != "" {
		q.Set("chatId", params.ChatID)
	}
	if params.MaxAttempts > 0 {
		q.Set("maxAttempts", strconv.Itoa(params.MaxAttempts))
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/messages/media/" + url.PathEscape(waMessageKey),
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data Media `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// Pin a message. params.Duration (seconds) is optional; params.ChatID is an
// optional locator hint.
func (r *MessagesResource) Pin(ctx context.Context, waMessageKey string, params PinMessageParams) (*OkResult, error) {
	q := url.Values{}
	if params.ChatID != "" {
		q.Set("chatId", params.ChatID)
	}
	var reqBody interface{}
	if params.Duration != nil {
		reqBody = map[string]int{"duration": *params.Duration}
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
	return decodeOkResult(body)
}

// Unpin a message. params.ChatID is an optional locator hint.
func (r *MessagesResource) Unpin(ctx context.Context, waMessageKey string, params UnpinMessageParams) (*OkResult, error) {
	q := url.Values{}
	if params.ChatID != "" {
		q.Set("chatId", params.ChatID)
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/messages/unpin/" + url.PathEscape(waMessageKey),
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	return decodeOkResult(body)
}

// ListPinned returns the pinned messages of a chat. Offset-paginated.
func (r *MessagesResource) ListPinned(ctx context.Context, chatID string) (*PinnedMessagesPage, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/messages/pinned/" + url.PathEscape(chatID),
	})
	if err != nil {
		return nil, err
	}
	out := &PinnedMessagesPage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// React adds or replaces a reaction on a message. Pass an empty Emoji to clear
// an existing reaction. params.ChatID is an optional locator hint.
func (r *MessagesResource) React(ctx context.Context, waMessageKey string, params ReactParams) (*OkResult, error) {
	q := url.Values{}
	if params.ChatID != "" {
		q.Set("chatId", params.ChatID)
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/messages/reactions/" + url.PathEscape(waMessageKey),
		query:  q,
		body:   map[string]string{"emoji": params.Emoji},
	})
	if err != nil {
		return nil, err
	}
	return decodeOkResult(body)
}

func decodeOkResult(body []byte) (*OkResult, error) {
	var env struct {
		Data OkResult `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}
