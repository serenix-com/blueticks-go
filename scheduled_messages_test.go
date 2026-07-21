package blueticks

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// smData is the canonical MessageResponse (SentMessage) fixture body.
func smData(id, msgType string) string {
	return `{"id":"` + id + `","waMessageKey":null,"to":"c1@c.us","type":"` + msgType + `","text":"hi","mediaUrl":null,"mediaKind":null,"pollQuestion":null,"pollOptions":null,"pollAllowMultiple":null,"status":"pending","sendAt":null,"createdAt":"2026-01-01T00:00:00Z","confirmedAt":null,"receivedAt":null,"readAt":null,"playedAt":null,"failedAt":null,"failureReason":null,"secret":null,"linkPreview":null}`
}

func TestScheduledMessagesCreate(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/scheduled-messages/c1@c.us" {
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
		if _, present := sent["chatId"]; present {
			t.Errorf("body must not include chatId (derived from path)")
		}
		w.WriteHeader(http.StatusCreated)
		jsonResponse(w, `{"success":true,"data":`+smData("sm_1", "text")+`}`)
	})
	got, err := c.ScheduledMessages.Create(context.Background(), "c1@c.us", SendMessageParams{
		Type: MessageTypeText, Text: "hi", IdempotencyKey: "key_1",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.ID != "sm_1" || got.Type != MessageTypeText || got.Status != MessageStatusPending {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestScheduledMessagesList(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/scheduled-messages" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("status"); got != "pending" {
			t.Errorf("status query: %q", got)
		}
		jsonResponse(w, `{"success":true,"data":[`+smData("sm_1", "text")+`],"limit":50,"skip":0,"total":1}`)
	})
	got, err := c.ScheduledMessages.List(context.Background(), ListScheduledMessagesParams{Limit: 5, Status: "pending"})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got.Data) != 1 || got.Data[0].ID != "sm_1" || got.Total != 1 {
		t.Fatalf("unexpected page: %+v", got)
	}
}

func TestScheduledMessagesRetrieve(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/scheduled-messages/sm_1" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"success":true,"data":`+smData("sm_1", "text")+`}`)
	})
	got, err := c.ScheduledMessages.Retrieve(context.Background(), "sm_1")
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
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
		jsonResponse(w, `{"success":true,"data":`+smData("sm_1", "text")+`}`)
	})
	got, err := c.ScheduledMessages.Update(context.Background(), "sm_1", UpdateScheduledMessageParams{Text: strptr("edited")})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if got.ID != "sm_1" {
		t.Fatalf("unexpected id: %q", got.ID)
	}
}

func TestScheduledMessagesDelete(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/v1/scheduled-messages/sm_1" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"success":true,"data":{"id":"sm_1","deleted":true}}`)
	})
	del, err := c.ScheduledMessages.Delete(context.Background(), "sm_1")
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if del.ID != "sm_1" || !del.Deleted {
		t.Fatalf("unexpected delete ref: %+v", del)
	}
}

func TestScheduledMessagesDecodeEmpty(t *testing.T) {
	if err := decodeJSON([]byte(`{}`), &SentMessage{}); err != nil {
		t.Fatalf("SentMessage empty decode: %v", err)
	}
}

func TestScheduledMessagesCreateUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.ScheduledMessages.Create(context.Background(), "c1@c.us", SendMessageParams{Type: MessageTypeText, Text: "hi"})
	expectAuthError(t, err)
}
