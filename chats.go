package blueticks

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// ChatsResource exposes /v1/chats and its chat-level write endpoints. Message
// operations live on MessagesResource.
type ChatsResource struct {
	client *Client
}

// ListChatsParams is the query-param bundle for ChatsResource.List.
type ListChatsParams struct {
	SearchToken         string
	Filter              string // "groups" | "contacts" | "newsletters"
	IncludeLastMessage  *bool
	IncludeExtendedInfo *bool
	IncludeWithoutName  *bool
	IncludeArchive      *bool
	Skip                int
	Limit               int
}

// ListParticipantsParams is the query-param bundle for ChatsResource.ListParticipants.
type ListParticipantsParams struct {
	SearchToken string
	Skip        int
	Limit       int
}

// List/search chats. Offset-paginated.
func (r *ChatsResource) List(ctx context.Context, params ListChatsParams) (*ChatsPage, error) {
	q := url.Values{}
	if params.SearchToken != "" {
		q.Set("searchToken", params.SearchToken)
	}
	if params.Filter != "" {
		q.Set("filter", params.Filter)
	}
	if params.IncludeLastMessage != nil {
		q.Set("includeLastMessage", strconv.FormatBool(*params.IncludeLastMessage))
	}
	if params.IncludeExtendedInfo != nil {
		q.Set("includeExtendedInfo", strconv.FormatBool(*params.IncludeExtendedInfo))
	}
	if params.IncludeWithoutName != nil {
		q.Set("includeWithoutName", strconv.FormatBool(*params.IncludeWithoutName))
	}
	if params.IncludeArchive != nil {
		q.Set("includeArchive", strconv.FormatBool(*params.IncludeArchive))
	}
	if params.Skip > 0 {
		q.Set("skip", strconv.Itoa(params.Skip))
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
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
	var env struct {
		Data Chat `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// Archive moves the chat to the archived list.
func (r *ChatsResource) Archive(ctx context.Context, chatID string) (*OkResult, error) {
	return r.okAction(ctx, chatID, "archive")
}

// Unarchive restores an archived chat to the main list.
func (r *ChatsResource) Unarchive(ctx context.Context, chatID string) (*OkResult, error) {
	return r.okAction(ctx, chatID, "unarchive")
}

// MarkRead clears the unread badge on the connected engine for the given chat.
func (r *ChatsResource) MarkRead(ctx context.Context, chatID string) (*OkResult, error) {
	return r.okAction(ctx, chatID, "mark_read")
}

func (r *ChatsResource) okAction(ctx context.Context, chatID, action string) (*OkResult, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/chats/" + url.PathEscape(chatID) + "/" + action,
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data OkResult `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// ListParticipants returns the participants of a chat. For group chats the list
// is paginated; for DMs the single counterparty is returned. Offset-paginated.
func (r *ChatsResource) ListParticipants(ctx context.Context, chatID string, params ListParticipantsParams) (*ParticipantsPage, error) {
	q := url.Values{}
	if params.SearchToken != "" {
		q.Set("searchToken", params.SearchToken)
	}
	if params.Skip > 0 {
		q.Set("skip", strconv.Itoa(params.Skip))
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
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
