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
	SearchToken    string
	IncludeArchive *bool
	Skip           int
	Limit          int
}

// List contacts, optionally filtered by a search token. Offset-paginated.
func (r *ContactsResource) List(ctx context.Context, params ListContactsParams) (*ContactsPage, error) {
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

// CommonGroupsParams are the optional query params for
// ContactsResource.CommonGroupsWithParams.
type CommonGroupsParams struct {
	Skip  int
	Limit int
}

// CommonGroups returns the groups the connected account shares with the given
// contact. Offset-paginated. Use CommonGroupsWithParams to page.
func (r *ContactsResource) CommonGroups(ctx context.Context, id string) (*CommonGroupsPage, error) {
	return r.CommonGroupsWithParams(ctx, id, CommonGroupsParams{})
}

// CommonGroupsWithParams is CommonGroups with offset pagination (Skip, Limit).
func (r *ContactsResource) CommonGroupsWithParams(ctx context.Context, id string, params CommonGroupsParams) (*CommonGroupsPage, error) {
	q := url.Values{}
	if params.Skip > 0 {
		q.Set("skip", strconv.Itoa(params.Skip))
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/contacts/" + url.PathEscape(id) + "/common_groups",
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	out := &CommonGroupsPage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}
