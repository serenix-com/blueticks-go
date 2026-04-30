package blueticks

import (
	"context"
	"net/http"
	"net/url"
)

// UtilsResource exposes /v1/utils.
type UtilsResource struct {
	client *Client
}

// ValidatePhone checks whether the given phone number or chat ID corresponds
// to a registered WhatsApp account and returns its canonical chat ID form.
func (r *UtilsResource) ValidatePhone(ctx context.Context, phoneOrChatID string) (*PhoneValidation, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/utils/validate_phone",
		body:   map[string]string{"phone_or_chat_id": phoneOrChatID},
	})
	if err != nil {
		return nil, err
	}
	out := &PhoneValidation{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// LinkPreview fetches OG/Twitter card metadata for an arbitrary URL — used
// by clients that want to render a richer preview inline before sending.
func (r *UtilsResource) LinkPreview(ctx context.Context, link string) (*LinkPreview, error) {
	q := url.Values{}
	q.Set("url", link)
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/utils/link_preview",
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	out := &LinkPreview{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}
