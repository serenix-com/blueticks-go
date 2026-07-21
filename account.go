package blueticks

import (
	"context"
	"net/http"
)

// AccountResource exposes /v1/account.
type AccountResource struct {
	client *Client
}

// Retrieve returns the account associated with the API key used for the request.
func (r *AccountResource) Retrieve(ctx context.Context) (*Account, error) {
	body, err := r.client.do(ctx, requestOptions{method: http.MethodGet, path: "/v1/account"})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data Account `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}
