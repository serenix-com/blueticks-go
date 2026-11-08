package blueticks

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// EnginesResource exposes /v1/engines.
type EnginesResource struct {
	client *Client
}

// List returns the WhatsApp engines connected to this workspace. The slice is
// empty when no engine is paired.
func (r *EnginesResource) List(ctx context.Context) ([]EngineStatus, error) {
	body, err := r.client.do(ctx, requestOptions{method: http.MethodGet, path: "/v1/engines"})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data []EngineStatus `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return env.Data, nil
}

// addQuery writes the non-empty engine-targeting fields onto q.
func (t EngineTarget) addQuery(q url.Values) {
	if t.Engine != "" {
		q.Set("engine", t.Engine)
	}
	if t.EngineID != "" {
		q.Set("engineId", t.EngineID)
	}
	if t.NumberID != "" {
		q.Set("numberId", t.NumberID)
	}
	if t.Phone != "" {
		q.Set("phone", t.Phone)
	}
	if t.FanoutLimit > 0 {
		q.Set("fanoutLimit", strconv.Itoa(t.FanoutLimit))
	}
}
