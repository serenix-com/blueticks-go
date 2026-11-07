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

// ListGroupsParams is the query-param bundle for GroupsResource.List.
type ListGroupsParams struct {
	SearchToken    string
	IncludeArchive *bool
	Skip           int
	Limit          int
}

// List the groups the connected WhatsApp engine sees. Offset-paginated.
func (r *GroupsResource) List(ctx context.Context, params ListGroupsParams) (*GroupsPage, error) {
	q := url.Values{}
	if params.SearchToken != "" {
		q.Set("searchToken", params.SearchToken)
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
	return decodeGroup(body)
}

// Get retrieves a group by id. Set params.Include="participants" to inline the
// members array.
func (r *GroupsResource) Get(ctx context.Context, groupID string, params GetGroupParams) (*Group, error) {
	q := url.Values{}
	if params.Include != "" {
		q.Set("include", params.Include)
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/groups/" + url.PathEscape(groupID),
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	return decodeGroup(body)
}

// Update mutates a group's name and/or settings.
func (r *GroupsResource) Update(ctx context.Context, groupID string, params UpdateGroupParams) (*Group, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPatch,
		path:   "/v1/groups/" + url.PathEscape(groupID),
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	return decodeGroup(body)
}

// AddMember adds one or more participants to the group. Returns the updated group.
func (r *GroupsResource) AddMember(ctx context.Context, groupID string, params AddMemberParams) (*Group, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/members",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	return decodeGroup(body)
}

// RemoveMember removes a participant from the group. Returns the updated group.
func (r *GroupsResource) RemoveMember(ctx context.Context, groupID, chatID string) (*Group, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/members/" + url.PathEscape(chatID),
	})
	if err != nil {
		return nil, err
	}
	return decodeGroup(body)
}

// PromoteAdmin grants admin privileges to a group member. Returns the updated group.
func (r *GroupsResource) PromoteAdmin(ctx context.Context, groupID, chatID string) (*Group, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/members/" + url.PathEscape(chatID) + "/admin",
	})
	if err != nil {
		return nil, err
	}
	return decodeGroup(body)
}

// DemoteAdmin revokes admin privileges from a group member. Returns the updated group.
func (r *GroupsResource) DemoteAdmin(ctx context.Context, groupID, chatID string) (*Group, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/members/" + url.PathEscape(chatID) + "/admin",
	})
	if err != nil {
		return nil, err
	}
	return decodeGroup(body)
}

// SetPicture replaces the group picture. Returns the updated group.
func (r *GroupsResource) SetPicture(ctx context.Context, groupID string, params SetGroupPictureParams) (*Group, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPut,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/picture",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	return decodeGroup(body)
}

// GetPicture returns the current group picture as a short-lived WhatsApp CDN
// URL. The result's URL is nil when the group has no picture. Requires
// groups:read.
func (r *GroupsResource) GetPicture(ctx context.Context, groupID string) (*GroupPicture, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/picture",
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data GroupPicture `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// RemovePicture removes the group picture. Gateway (Baileys) engine only: a
// Chrome-extension-only account gets 501 unsupported_command. Only group admins
// can when the group is admins-only-edit (403 not_authorized). Requires
// groups:write.
func (r *GroupsResource) RemovePicture(ctx context.Context, groupID string) (*GroupPictureRemoved, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/picture",
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data GroupPictureRemoved `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// Leave removes the connected identity from the group. Idempotent — the
// endpoint returns 204 with no body.
func (r *GroupsResource) Leave(ctx context.Context, groupID string) error {
	_, err := r.client.do(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/v1/groups/" + url.PathEscape(groupID) + "/members/me",
	})
	return err
}

func decodeGroup(body []byte) (*Group, error) {
	var env struct {
		Data Group `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}
