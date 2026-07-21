package blueticks

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// NewslettersResource exposes /v1/newsletters.
type NewslettersResource struct {
	client *Client
}

// ListNewslettersParams is the query-param bundle for NewslettersResource.List.
type ListNewslettersParams struct {
	SearchToken    string
	IncludeArchive *bool
	Skip           int
	Limit          int
}

// Create a new WhatsApp newsletter (channel) via the connected engine.
func (r *NewslettersResource) Create(ctx context.Context, params CreateNewsletterParams) (*Newsletter, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/newsletters",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data Newsletter `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// List newsletters visible to the connected WhatsApp engine. Offset-paginated.
func (r *NewslettersResource) List(ctx context.Context, params ListNewslettersParams) (*NewslettersPage, error) {
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
		path:   "/v1/newsletters",
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	out := &NewslettersPage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Retrieve returns a newsletter by its JID.
func (r *NewslettersResource) Retrieve(ctx context.Context, id string) (*Newsletter, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/newsletters/" + url.PathEscape(id),
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data Newsletter `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}
