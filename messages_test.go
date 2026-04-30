package blueticks

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestMessagesSend(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/messages" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Idempotency-Key"); got != "key_1" {
			t.Errorf("missing idempotency-key: %q", got)
		}
		body, _ := io.ReadAll(r.Body)
		var sent map[string]any
		_ = json.Unmarshal(body, &sent)
		if sent["to"] != "+15555550100" || sent["text"] != "hi" {
			t.Errorf("unexpected body: %v", sent)
		}
		jsonResponse(w, `{"id":"msg_1","to":"+15555550100","text":"hi","status":"queued","created_at":"2026-01-01T00:00:00Z"}`)
	})
	got, err := c.Messages.Send(context.Background(), SendMessageParams{
		To: "+15555550100", Text: "hi", IdempotencyKey: "key_1",
	})
	if err != nil {
		t.Fatalf("Send: %v", err)
	}
	if got.ID != "msg_1" || got.Status != "queued" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestMessagesGet(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/messages/msg_1" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"id":"msg_1","to":"+1","status":"delivered","created_at":"2026-01-01T00:00:00Z"}`)
	})
	got, err := c.Messages.Get(context.Background(), "msg_1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.ID != "msg_1" {
		t.Fatalf("unexpected id: %q", got.ID)
	}
}

func TestMessagesList(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/messages" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "5" {
			t.Errorf("limit query: %q", got)
		}
		jsonResponse(w, `{"data":[{"id":"msg_1","to":"+1","status":"delivered","created_at":"2026-01-01T00:00:00Z"}],"has_more":false}`)
	})
	got, err := c.Messages.List(context.Background(), ListMessagesParams{Limit: 5})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got.Data) != 1 || got.HasMore {
		t.Fatalf("unexpected page: %+v", got)
	}
}

func TestMessagesUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Messages.Send(context.Background(), SendMessageParams{To: "+1", Text: "hi"})
	expectAuthError(t, err)
}
