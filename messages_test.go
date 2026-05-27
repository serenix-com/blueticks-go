package blueticks

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// msgFixture is the canonical Message fixture matching the spec's 200/201 response.
const msgFixture = `{"id":"msg_1","key":null,"to":"+15555550100","from":null,"type":"text","text":"hi","media_url":null,"media_kind":null,"poll_question":null,"status":"queued","send_at":null,"created_at":"2026-01-01T00:00:00Z","sent_at":null,"delivered_at":null,"read_at":null,"failed_at":null,"failure_reason":null,"link_preview":null}`

func TestMessagesSendText(t *testing.T) {
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
		if sent["type"] != "text" || sent["to"] != "+15555550100" || sent["text"] != "hi" {
			t.Errorf("unexpected body: %v", sent)
		}
		jsonResponse(w, msgFixture)
	})
	got, err := c.Messages.Send(context.Background(), SendMessageParams{
		Type: MessageTypeText, To: "+15555550100", Text: "hi", IdempotencyKey: "key_1",
	})
	if err != nil {
		t.Fatalf("Send: %v", err)
	}
	if got.ID != "msg_1" || got.Status != "queued" || got.Type != MessageTypeText {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestMessagesSendMedia(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/messages" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var sent map[string]any
		_ = json.Unmarshal(body, &sent)
		if sent["type"] != "media" {
			t.Errorf("expected type=media, got %v", sent["type"])
		}
		mediaFix := `{"id":"msg_2","key":null,"to":"+15555550100","from":null,"type":"media","text":null,"media_url":"https://cdn.example.com/x.pdf","media_kind":"document","poll_question":null,"status":"queued","send_at":null,"created_at":"2026-01-01T00:00:00Z","sent_at":null,"delivered_at":null,"read_at":null,"failed_at":null,"failure_reason":null,"link_preview":null}`
		jsonResponse(w, mediaFix)
	})
	got, err := c.Messages.Send(context.Background(), SendMessageParams{
		Type:  MessageTypeMedia,
		To:    "+15555550100",
		Media: &SendMediaObject{URL: "https://cdn.example.com/x.pdf", Kind: "document", Filename: "receipt.pdf"},
	})
	if err != nil {
		t.Fatalf("Send media: %v", err)
	}
	if got.ID != "msg_2" || got.Type != MessageTypeMedia {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestMessagesSendPoll(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/messages" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var sent map[string]any
		_ = json.Unmarshal(body, &sent)
		if sent["type"] != "poll" {
			t.Errorf("expected type=poll, got %v", sent["type"])
		}
		pollFix := `{"id":"msg_3","key":null,"to":"+15555550100","from":null,"type":"poll","text":null,"media_url":null,"media_kind":null,"poll_question":"Pizza?","status":"queued","send_at":null,"created_at":"2026-01-01T00:00:00Z","sent_at":null,"delivered_at":null,"read_at":null,"failed_at":null,"failure_reason":null,"link_preview":null}`
		jsonResponse(w, pollFix)
	})
	got, err := c.Messages.Send(context.Background(), SendMessageParams{
		Type: MessageTypePoll,
		To:   "+15555550100",
		Poll: &SendPollObject{Question: "Pizza?", Options: []string{"Yes", "No"}},
	})
	if err != nil {
		t.Fatalf("Send poll: %v", err)
	}
	if got.ID != "msg_3" || got.Type != MessageTypePoll {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestMessagesGet(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/messages/msg_1" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, msgFixture)
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
		jsonResponse(w, `{"data":[`+msgFixture+`],"has_more":false}`)
	})
	got, err := c.Messages.List(context.Background(), ListMessagesParams{Limit: 5})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got.Data) != 1 || got.HasMore {
		t.Fatalf("unexpected page: %+v", got)
	}
}

func TestMessagesDecodeEmptyResponse(t *testing.T) {
	if err := decodeJSON([]byte(`{}`), &Message{}); err != nil {
		t.Fatalf("Message empty decode: %v", err)
	}
}

func TestMessagesUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Messages.Send(context.Background(), SendMessageParams{Type: MessageTypeText, To: "+1", Text: "hi"})
	expectAuthError(t, err)
}

func TestMessagesUpdate(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/v1/messages/msg_1" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, msgFixture)
	})
	newText := "edited"
	got, err := c.Messages.Update(context.Background(), "msg_1", UpdateMessageParams{Text: &newText})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if got.ID != "msg_1" {
		t.Fatalf("unexpected id: %q", got.ID)
	}
}

func TestMessagesUpdateUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	newText := "edited"
	_, err := c.Messages.Update(context.Background(), "msg_1", UpdateMessageParams{Text: &newText})
	expectAuthError(t, err)
}
