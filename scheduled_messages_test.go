package blueticks

import (
	"context"
	"net/http"
	"testing"
)

func TestScheduledMessagesList(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/scheduled-messages" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "5" {
			t.Errorf("limit query: %q", got)
		}
		jsonResponse(w, `{"data":[{"id":"sm_1","to":"+1","text":"hi","media_url":null,"media_caption":null,"media_filename":null,"media_mime_type":null,"send_at":"2026-01-02T00:00:00Z","status":"scheduled","is_recurring":false,"recurrence_rule":null,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}],"has_more":false,"next_cursor":null}`)
	})
	got, err := c.ScheduledMessages.List(context.Background(), ListScheduledMessagesParams{Limit: 5})
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
		jsonResponse(w, `{"id":"sm_1","to":"+1","text":"hi","media_url":null,"media_caption":null,"media_filename":null,"media_mime_type":null,"send_at":"2026-01-02T00:00:00Z","status":"scheduled","is_recurring":false,"recurrence_rule":null,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}`)
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
		jsonResponse(w, `{"id":"sm_1","to":"+1","text":"updated","media_url":null,"media_caption":null,"media_filename":null,"media_mime_type":null,"send_at":"2026-02-02T00:00:00Z","status":"scheduled","is_recurring":false,"recurrence_rule":null,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-02T00:00:00Z"}`)
	})
	newText := "updated"
	got, err := c.ScheduledMessages.Update(context.Background(), "sm_1", UpdateScheduledMessageParams{Text: &newText})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if got.Text == nil || *got.Text != "updated" {
		t.Fatalf("unexpected text: %+v", got.Text)
	}
}

func TestScheduledMessagesDelete(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/v1/scheduled-messages/sm_1" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"id":"sm_1","deleted":true}`)
	})
	got, err := c.ScheduledMessages.Delete(context.Background(), "sm_1")
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if got.ID != "sm_1" || !got.Deleted {
		t.Fatalf("unexpected delete ref: %+v", got)
	}
}

func TestScheduledMessagesDecodeEmpty(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, `{}`)
	})
	got, err := c.ScheduledMessages.Retrieve(context.Background(), "sm_1")
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	if got.ID != "" {
		t.Fatalf("expected zero ID, got %q", got.ID)
	}
}

func TestScheduledMessagesUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.ScheduledMessages.List(context.Background(), ListScheduledMessagesParams{})
	expectAuthError(t, err)
}
