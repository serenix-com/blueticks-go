package blueticks

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// AudiencesResource exposes /v1/audiences.
type AudiencesResource struct {
	client *Client
}

// ListAudiencesParams is the query-param bundle for AudiencesResource.List.
type ListAudiencesParams struct {
	Order string // "asc" or "desc" (default "desc")
	Skip  int
	Limit int
}

// Create an audience, optionally with seed contacts.
func (r *AudiencesResource) Create(ctx context.Context, params CreateAudienceParams) (*Audience, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/audiences",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data Audience `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// List audiences, offset-paginated.
func (r *AudiencesResource) List(ctx context.Context, params ListAudiencesParams) (*AudiencesPage, error) {
	q := url.Values{}
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
		path:   "/v1/audiences",
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	out := &AudiencesPage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get retrieves an audience by id.
func (r *AudiencesResource) Get(ctx context.Context, id string) (*Audience, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/audiences/" + url.PathEscape(id),
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data Audience `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// Update renames an audience.
func (r *AudiencesResource) Update(ctx context.Context, id string, params UpdateAudienceParams) (*Audience, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPatch,
		path:   "/v1/audiences/" + url.PathEscape(id),
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data Audience `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// Delete an audience. Returns the deleted ref ({id, deleted: true}).
func (r *AudiencesResource) Delete(ctx context.Context, id string) (*DeletedResource, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/v1/audiences/" + url.PathEscape(id),
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

// AppendContacts adds contacts to an existing audience and returns the number
// added plus the new total contact count.
func (r *AudiencesResource) AppendContacts(ctx context.Context, id string, contacts []ContactInput) (*AppendContactsResult, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/audiences/" + url.PathEscape(id) + "/contacts",
		body:   map[string]interface{}{"contacts": contacts},
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data AppendContactsResult `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// UpdateContact updates a single contact within an audience.
func (r *AudiencesResource) UpdateContact(ctx context.Context, audienceID, contactID string, params UpdateAudienceContactParams) (*AudienceContact, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPatch,
		path:   "/v1/audiences/" + url.PathEscape(audienceID) + "/contacts/" + url.PathEscape(contactID),
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data AudienceContact `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// DeleteContact removes a contact from an audience. Returns the deleted ref
// ({id}).
func (r *AudiencesResource) DeleteContact(ctx context.Context, audienceID, contactID string) (*DeletedRef, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/v1/audiences/" + url.PathEscape(audienceID) + "/contacts/" + url.PathEscape(contactID),
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data DeletedRef `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}
