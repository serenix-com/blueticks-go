package blueticks

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// GroupsResource exposes /v1/groups.
type GroupsResource struct {
	client *Client
}

// AddMemberParams is the body for GroupsResource.AddMember.
type AddMemberParams struct {
	ChatID string `json:"chat_id"`
}

// ListGroupsParams is the query-param bundle for GroupsResource.List.
type ListGroupsParams struct {
	Limit  int
	Cursor string
	Q      string // case-insensitive substring match against group name
}

// SetGroupPictureParams is the body for GroupsResource.SetPicture.
//
// FileDataURL is required and must be a `data:<mime>;base64,<data>` URL.
type SetGroupPictureParams struct {
	FileDataURL  string `json:"file_data_url"`
	FileName     string `json:"file_name,omitempty"`
	FileMimeType string `json:"file_mime_type,omitempty"`
}

// List the groups the connected WhatsApp engine sees.
//
// Supports cursor pagination (Limit + Cursor) and an optional case-insensitive
// substring search on the group name via Q.
func (r *GroupsResource) List(ctx context.Context, params ListGroupsParams) (*GroupsPage, error) {
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		q.Set("cursor", params.Cursor)
	}
	if params.Q != "" {
		q.Set("q", params.Q)
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/groups",
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	out := &GroupsPage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Create a new WhatsApp group with the given participants.
func (r *GroupsResource) Create(ctx context.Context, params CreateGroupParams) (*Group, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/groups",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	out := &Group{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get retrieves a group by id.
func (r *GroupsResource) Get(ctx context.Context, groupID string) (*Group, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/groups/" + url.PathEscape(groupID),
	})
	if err != nil {
		return nil, err
	}
	out := &Group{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Update mutates a group's name and/or settings. Settings is a map of
// flags (announce, restrict, …).
func (r *GroupsResource) Update(ctx context.Context, groupID string, params UpdateGroupParams) (*Group, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPatch,
		path:   "/v1/groups/" + url.PathEscape(groupID),
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	out := &Group{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// AddMember adds a participant to the group by chat_id (JID) or phone number in international format (e.g. +14155551234).
// Returns the updated group.
func (r *GroupsResource) AddMember(ctx context.Context, groupID string, params AddMemberParams) (*Group, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/members",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	out := &Group{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// RemoveMember removes a participant from the group. Returns the updated
// group.
func (r *GroupsResource) RemoveMember(ctx context.Context, groupID, chatID string) (*Group, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/members/" + url.PathEscape(chatID),
	})
	if err != nil {
		return nil, err
	}
	out := &Group{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// PromoteAdmin grants admin privileges to a group member. Returns the
// updated group.
func (r *GroupsResource) PromoteAdmin(ctx context.Context, groupID, chatID string) (*Group, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/members/" + url.PathEscape(chatID) + "/admin",
	})
	if err != nil {
		return nil, err
	}
	out := &Group{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// DemoteAdmin revokes admin privileges from a group member. Returns the
// updated group.
func (r *GroupsResource) DemoteAdmin(ctx context.Context, groupID, chatID string) (*Group, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/members/" + url.PathEscape(chatID) + "/admin",
	})
	if err != nil {
		return nil, err
	}
	out := &Group{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// SetPicture replaces the group picture. FileDataURL must be a base64 data
// URL (PNG/JPEG, ≤20 MiB). Returns the updated group.
func (r *GroupsResource) SetPicture(ctx context.Context, groupID string, params SetGroupPictureParams) (*Group, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPut,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/picture",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	out := &Group{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Leave removes the authenticated identity from the group. Idempotent —
// succeeds even if already not a member. The endpoint returns 204 with no
// body.
func (r *GroupsResource) Leave(ctx context.Context, groupID string) error {
	_, err := r.client.do(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/members/me",
	})
	return err
}
