package blueticks

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// smFixture is the canonical ScheduledMessage fixture matching the spec's
// 200/201 response.
const smFixture = `{"id":"sm_1","key":null,"to":"+15555550100","from":null,"type":"text","text":"hi","mediaUrl":null,"mediaKind":null,"pollQuestion":null,"status":"pending","sendAt":null,"createdAt":"2026-01-01T00:00:00Z","confirmedAt":null,"receivedAt":null,"readAt":null,"playedAt":null,"failedAt":null,"failureReason":null,"linkPreview":null}`

func TestScheduledMessagesCreateText(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/scheduled-messages/+15555550100" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Idempotency-Key"); got != "key_1" {
			t.Errorf("missing idempotency-key: %q", got)
		}
		body, _ := io.ReadAll(r.Body)
		var sent map[string]any
		_ = json.Unmarshal(body, &sent)
		if sent["type"] != "text" || sent["text"] != "hi" {
			t.Errorf("unexpected body: %v", sent)
		}
		jsonResponse(w, smFixture)
	})
	got, err := c.ScheduledMessages.Create(context.Background(), "+15555550100", CreateScheduledMessageParams{
		Type: MessageTypeText, Text: "hi", IdempotencyKey: "key_1",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.ID != "sm_1" || got.Status != "pending" || got.Type != MessageTypeText {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestScheduledMessagesCreateMedia(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/scheduled-messages/+15555550100" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var sent map[string]any
		_ = json.Unmarshal(body, &sent)
		if sent["type"] != "media" {
			t.Errorf("expected type=media, got %v", sent["type"])
		}
		mediaFix := `{"id":"sm_2","key":null,"to":"+15555550100","from":null,"type":"media","text":null,"mediaUrl":"https://cdn.example.com/x.pdf","mediaKind":"document","pollQuestion":null,"status":"pending","sendAt":null,"createdAt":"2026-01-01T00:00:00Z","confirmedAt":null,"receivedAt":null,"readAt":null,"playedAt":null,"failedAt":null,"failureReason":null,"linkPreview":null}`
		jsonResponse(w, mediaFix)
	})
	got, err := c.ScheduledMessages.Create(context.Background(), "+15555550100", CreateScheduledMessageParams{
		Type:  MessageTypeMedia,
		Media: &SendMediaObject{URL: "https://cdn.example.com/x.pdf", Kind: "document", Filename: "receipt.pdf"},
	})
	if err != nil {
		t.Fatalf("Create media: %v", err)
	}
	if got.ID != "sm_2" || got.Type != MessageTypeMedia {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestScheduledMessagesCreatePoll(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/scheduled-messages/+15555550100" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var sent map[string]any
		_ = json.Unmarshal(body, &sent)
		if sent["type"] != "poll" {
			t.Errorf("expected type=poll, got %v", sent["type"])
		}
		pollFix := `{"id":"sm_3","key":null,"to":"+15555550100","from":null,"type":"poll","text":null,"mediaUrl":null,"mediaKind":null,"pollQuestion":"Pizza?","status":"pending","sendAt":null,"createdAt":"2026-01-01T00:00:00Z","confirmedAt":null,"receivedAt":null,"readAt":null,"playedAt":null,"failedAt":null,"failureReason":null,"linkPreview":null}`
		jsonResponse(w, pollFix)
	})
	got, err := c.ScheduledMessages.Create(context.Background(), "+15555550100", CreateScheduledMessageParams{
		Type: MessageTypePoll,
		Poll: &SendPollObject{Question: "Pizza?", Options: []string{"Yes", "No"}},
	})
	if err != nil {
		t.Fatalf("Create poll: %v", err)
	}
	if got.ID != "sm_3" || got.Type != MessageTypePoll {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestScheduledMessagesList(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/scheduled-messages" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "5" {
			t.Errorf("limit query: %q", got)
		}
		if got := r.URL.Query().Get("status"); got != "pending" {
			t.Errorf("status query: %q", got)
		}
		jsonResponse(w, `{"data":[`+smFixture+`],"hasMore":false}`)
	})
	got, err := c.ScheduledMessages.List(context.Background(), ListScheduledMessagesParams{Limit: 5, Status: "pending"})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got.Data) != 1 || got.HasMore {
		t.Fatalf("unexpected page: %+v", got)
	}
	if got.Data[0].ID != "sm_1" {
		t.Fatalf("unexpected id: %q", got.Data[0].ID)
	}
}

func TestScheduledMessagesRetrieve(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/scheduled-messages/sm_1" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, smFixture)
	})
	got, err := c.ScheduledMessages.Retrieve(context.Background(), "sm_1")
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	if got.ID != "sm_1" {
		t.Fatalf("unexpected id: %q", got.ID)
	}
	if got.Text == nil || *got.Text != "hi" {
		t.Fatalf("unexpected text: %+v", got.Text)
	}
}

func TestScheduledMessagesUpdate(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/v1/scheduled-messages/sm_1" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, smFixture)
	})
	newText := "edited"
	got, err := c.ScheduledMessages.Update(context.Background(), "sm_1", UpdateScheduledMessageParams{Text: &newText})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if got.ID != "sm_1" {
		t.Fatalf("unexpected id: %q", got.ID)
	}
}

func TestScheduledMessagesDecodeEmpty(t *testing.T) {
	if err := decodeJSON([]byte(`{}`), &ScheduledMessage{}); err != nil {
		t.Fatalf("ScheduledMessage empty decode: %v", err)
	}
}

func TestScheduledMessagesCreateUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.ScheduledMessages.Create(context.Background(), "+1", CreateScheduledMessageParams{Type: MessageTypeText, Text: "hi"})
	expectAuthError(t, err)
}

func TestScheduledMessagesUpdateUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	newText := "edited"
	_, err := c.ScheduledMessages.Update(context.Background(), "sm_1", UpdateScheduledMessageParams{Text: &newText})
	expectAuthError(t, err)
}
