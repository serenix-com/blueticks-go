package blueticks

import (
	"context"
	"net/http"
	"net/url"
)

// GroupsResource exposes /v1/groups.
type GroupsResource struct {
	client *Client
}

// SetGroupPictureParams is the body for GroupsResource.SetPicture.
//
// FileDataURL is required and must be a `data:<mime>;base64,<data>` URL.
type SetGroupPictureParams struct {
	FileDataURL  string `json:"file_data_url"`
	FileName     string `json:"file_name,omitempty"`
	FileMimeType string `json:"file_mime_type,omitempty"`
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

// AddMember adds a participant to a group.
func (r *GroupsResource) AddMember(ctx context.Context, groupID, chatID string) (map[string]any, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/members",
		body:   map[string]string{"chat_id": chatID},
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

// RemoveMember removes a participant from a group.
func (r *GroupsResource) RemoveMember(ctx context.Context, groupID, chatID string) (map[string]any, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/members/" + url.PathEscape(chatID),
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

// PromoteAdmin grants group-admin to an existing member.
func (r *GroupsResource) PromoteAdmin(ctx context.Context, groupID, chatID string) (map[string]any, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/members/" + url.PathEscape(chatID) + "/admin",
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

// DemoteAdmin revokes group-admin from a member.
func (r *GroupsResource) DemoteAdmin(ctx context.Context, groupID, chatID string) (map[string]any, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/members/" + url.PathEscape(chatID) + "/admin",
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

// SetPicture replaces a group's avatar. FileDataURL is required.
func (r *GroupsResource) SetPicture(ctx context.Context, groupID string, params SetGroupPictureParams) (map[string]any, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPut,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/picture",
		body:   params,
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

// Leave removes the authenticated user from the group.
func (r *GroupsResource) Leave(ctx context.Context, groupID string) (map[string]any, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/members/me",
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
