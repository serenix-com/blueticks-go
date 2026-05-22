package blueticks

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// ContactsResource exposes /v1/contacts.
type ContactsResource struct {
	client *Client
}

// ListContactsParams is the query-param bundle for ContactsResource.List.
type ListContactsParams struct {
	Query  string
	Limit  int
	Cursor string
}

// List contacts, optionally filtered by free-text query. Cursor-paginated.
func (r *ContactsResource) List(ctx context.Context, params ListContactsParams) (*ContactsPage, error) {
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
		path:   "/v1/contacts",
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	out := &ContactsPage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}
