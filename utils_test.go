package blueticks

import (
	"context"
	"net/http"
	"testing"
)

func TestUtilsValidatePhone(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/utils/validate_phone" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"valid":true,"formatted_chat_id":"15555550100@c.us"}`)
	})
	got, err := c.Utils.ValidatePhone(context.Background(), "+15555550100")
	if err != nil {
		t.Fatalf("ValidatePhone: %v", err)
	}
	if !got.Valid || got.FormattedChatID == "" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestUtilsLinkPreview(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/utils/link_preview" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("url") != "https://example.com" {
			t.Errorf("expected url query")
		}
		jsonResponse(w, `{"title":"Example","description":"d","thumbnail":"t","canonical_url":"https://example.com"}`)
	})
	got, err := c.Utils.LinkPreview(context.Background(), "https://example.com")
	if err != nil {
		t.Fatalf("LinkPreview: %v", err)
	}
	if got.Title != "Example" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestUtilsUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Utils.ValidatePhone(context.Background(), "+1")
	expectAuthError(t, err)
}
