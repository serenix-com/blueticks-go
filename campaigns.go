package blueticks

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// CampaignsResource exposes /v1/campaigns.
type CampaignsResource struct {
	client *Client
}

// ListCampaignsParams is the query-param bundle for CampaignsResource.List.
type ListCampaignsParams struct {
	Limit  int
	Cursor string
}

// Create a campaign.
func (r *CampaignsResource) Create(ctx context.Context, params CreateCampaignParams) (*Campaign, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/campaigns",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	out := &Campaign{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// List campaigns, newest first. Cursor-paginated.
func (r *CampaignsResource) List(ctx context.Context, params ListCampaignsParams) (*CampaignsPage, error) {
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		q.Set("cursor", params.Cursor)
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/campaigns",
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	out := &CampaignsPage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get retrieves a campaign by id.
func (r *CampaignsResource) Get(ctx context.Context, id string) (*Campaign, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/campaigns/" + url.PathEscape(id),
	})
	if err != nil {
		return nil, err
	}
	out := &Campaign{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Pause a running campaign.
func (r *CampaignsResource) Pause(ctx context.Context, id string) (*Campaign, error) {
	return r.transition(ctx, id, "pause")
}

// Resume a paused campaign.
func (r *CampaignsResource) Resume(ctx context.Context, id string) (*Campaign, error) {
	return r.transition(ctx, id, "resume")
}

// Cancel an in-flight or paused campaign.
func (r *CampaignsResource) Cancel(ctx context.Context, id string) (*Campaign, error) {
	return r.transition(ctx, id, "cancel")
}

func (r *CampaignsResource) transition(ctx context.Context, id, action string) (*Campaign, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/campaigns/" + url.PathEscape(id) + "/" + action,
	})
	if err != nil {
		return nil, err
	}
	out := &Campaign{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}
