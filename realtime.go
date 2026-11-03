package blueticks

import (
	"context"
	"net/http"
)

// RealtimeResource exposes /v1/realtime — receiving messages over MQTT when no
// public webhook URL is available.
type RealtimeResource struct {
	client *Client
}

// Token trades the API key for an MQTT credential: broker URL, username,
// password and the topic to subscribe to. Your process dials out and
// subscribes, so nothing needs to be exposed. The credential does not
// practically expire, so fetching it once at start-up is enough.
func (r *RealtimeResource) Token(ctx context.Context) (*RealtimeToken, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/realtime/token",
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data RealtimeToken `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}
