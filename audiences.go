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
	Limit  int
	Cursor string
}

// GetAudienceParams provides optional pagination over the embedded contacts.
type GetAudienceParams struct {
	// Page is 1-based. Zero means "use server default".
	Page int
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
	out := &Audience{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// List audiences, newest first. Cursor-paginated.
func (r *AudiencesResource) List(ctx context.Context, params ListAudiencesParams) (*AudiencesPage, error) {
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		q.Set("cursor", params.Cursor)
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

// Get retrieves an audience by id, optionally returning a specific page of
// embedded contacts.
func (r *AudiencesResource) Get(ctx context.Context, id string, params GetAudienceParams) (*Audience, error) {
	q := url.Values{}
	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/audiences/" + url.PathEscape(id),
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	out := &Audience{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
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
	out := &Audience{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete an audience.
func (r *AudiencesResource) Delete(ctx context.Context, id string) error {
	_, err := r.client.do(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/v1/audiences/" + url.PathEscape(id),
	})
	return err
}

// AppendContacts adds contacts to an existing audience and returns the
// number added plus the new total contact count.
func (r *AudiencesResource) AppendContacts(ctx context.Context, id string, contacts []AudienceContactInput) (*AppendContactsResult, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/audiences/" + url.PathEscape(id) + "/contacts",
		body:   map[string]interface{}{"contacts": contacts},
	})
	if err != nil {
		return nil, err
	}
	out := &AppendContactsResult{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
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
	out := &AudienceContact{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteContact removes a contact from an audience.
func (r *AudiencesResource) DeleteContact(ctx context.Context, audienceID, contactID string) error {
	_, err := r.client.do(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/v1/audiences/" + url.PathEscape(audienceID) + "/contacts/" + url.PathEscape(contactID),
	})
	return err
}
