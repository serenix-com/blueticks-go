package blueticks

import (
	"context"
	"net/http"
)

// WaLabelsResource exposes /v1/wa-labels: the native WhatsApp Business labels
// defined on the connected account.
type WaLabelsResource struct {
	client *Client
}

// List returns the WhatsApp Business labels defined on the connected account. Pass
// their ids to ChatsResource.UpdateWaLabels or SendMessageParams.WaLabelIDs.
func (r *WaLabelsResource) List(ctx context.Context) ([]WaLabel, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/wa-labels",
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data []WaLabel `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return env.Data, nil
}
