package blueticks

import (
	"context"
	"net/http"
	"testing"
)

func TestWebhooksFlow(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/webhooks":
			jsonResponse(w, `{"id":"wh_1","url":"https://x","events":["message.delivered"],"status":"active","created_at":"2026-01-01T00:00:00Z","secret":"whsec_abc"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/webhooks":
			jsonResponse(w, `{"data":[{"id":"wh_1","url":"https://x","events":["message.delivered"],"status":"active","created_at":"2026-01-01T00:00:00Z"}],"has_more":false}`)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/webhooks/wh_1":
			jsonResponse(w, `{"id":"wh_1","url":"https://x","events":["message.delivered"],"status":"active","created_at":"2026-01-01T00:00:00Z"}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/v1/webhooks/wh_1":
			jsonResponse(w, `{"id":"wh_1","url":"https://y","events":["message.delivered"],"status":"active","created_at":"2026-01-01T00:00:00Z"}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/webhooks/wh_1":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/webhooks/wh_1/rotate-secret":
			jsonResponse(w, `{"id":"wh_1","url":"https://x","events":["message.delivered"],"status":"active","created_at":"2026-01-01T00:00:00Z","secret":"whsec_new"}`)
		default:
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	wh, err := c.Webhooks.Create(ctx, CreateWebhookParams{URL: "https://x", Events: []string{"message.delivered"}})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if wh.Secret != "whsec_abc" {
		t.Fatalf("expected secret, got %q", wh.Secret)
	}
	if _, err := c.Webhooks.List(ctx, ListWebhooksParams{}); err != nil {
		t.Fatalf("List: %v", err)
	}
	if _, err := c.Webhooks.Get(ctx, "wh_1"); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if _, err := c.Webhooks.Update(ctx, "wh_1", UpdateWebhookParams{URL: "https://y"}); err != nil {
		t.Fatalf("Update: %v", err)
	}
	if err := c.Webhooks.Delete(ctx, "wh_1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	rot, err := c.Webhooks.RotateSecret(ctx, "wh_1")
	if err != nil {
		t.Fatalf("RotateSecret: %v", err)
	}
	if rot.Secret != "whsec_new" {
		t.Fatalf("expected new secret, got %q", rot.Secret)
	}
}

func TestWebhooksUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Webhooks.List(context.Background(), ListWebhooksParams{})
	expectAuthError(t, err)
}
