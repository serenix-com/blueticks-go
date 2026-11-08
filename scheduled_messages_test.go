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

func TestScheduledMessagesCreateMentionsAndRecipientName(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		_ = json.Unmarshal(raw, &body)
		m, _ := body["mentions"].([]any)
		if len(m) != 1 || m[0] != "12025550101@c.us" || body["recipientName"] != "Team" {
			t.Errorf("unexpected body: %s", raw)
		}
		w.WriteHeader(http.StatusCreated)
		jsonResponse(w, `{"success":true,"data":{"id":"m1","to":"g1@g.us","type":"text","status":"pending","createdAt":"2026-01-01T00:00:00Z"}}`)
	})
	_, err := c.ScheduledMessages.Create(context.Background(), "g1@g.us", SendMessageParams{
		Type: MessageTypeText, Text: "hi @Jon", Mentions: []string{"12025550101@c.us"}, RecipientName: "Team",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
}

func TestScheduledMessagesUpdateMentions(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/v1/scheduled-messages/m1" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		_ = json.Unmarshal(raw, &body)
		if m, _ := body["mentions"].([]any); len(m) != 1 {
			t.Errorf("unexpected body: %s", raw)
		}
		jsonResponse(w, `{"success":true,"data":{"id":"m1","to":"g1@g.us","type":"text","status":"pending","createdAt":"2026-01-01T00:00:00Z"}}`)
	})
	if _, err := c.ScheduledMessages.Update(context.Background(), "m1", UpdateScheduledMessageParams{Mentions: []string{"12025550101@c.us"}}); err != nil {
		t.Fatalf("Update: %v", err)
	}
}

func TestScheduledMessagesCancel(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/scheduled-messages/sm_1/cancel" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"success":true,"data":{"id":"sm_1","waMessageKey":null,"to":"c1@c.us","type":"text","text":"hi","mediaUrl":null,"mediaKind":null,"pollQuestion":null,"pollOptions":null,"pollAllowMultiple":null,"status":"cancelled","sendAt":null,"createdAt":"2026-01-01T00:00:00Z","confirmedAt":null,"receivedAt":null,"readAt":null,"playedAt":null,"failedAt":null,"failureReason":null,"secret":null,"linkPreview":null}}`)
	})
	got, err := c.ScheduledMessages.Cancel(context.Background(), "sm_1")
	if err != nil {
		t.Fatalf("Cancel: %v", err)
	}
	if got.ID != "sm_1" || got.Status != MessageStatusCancelled {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestScheduledMessagesListIncludeSentAndPreview(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("includeSent") != "true" || r.URL.Query().Get("includePreview") != "true" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		jsonResponse(w, `{"success":true,"data":[],"limit":50,"skip":0,"total":0}`)
	})
	yes := true
	if _, err := c.ScheduledMessages.List(context.Background(), ListScheduledMessagesParams{IncludeSent: &yes, IncludePreview: &yes}); err != nil {
		t.Fatalf("List: %v", err)
	}
}

func TestScheduledMessagesCreateRecurrenceWaLabelsAndFrom(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/scheduled-messages/c1@c.us" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		var sent map[string]any
		_ = json.Unmarshal(raw, &sent)
		if _, present := sent["mediaFileId"]; present {
			t.Errorf("unset mediaFileId must be omitted: %s", raw)
		}
		if sent["from"] != "+15550001234" {
			t.Errorf("from not sent: %s", raw)
		}
		rec, _ := sent["recurrence"].(map[string]any)
		if rec["frequency"] != "weekly" || rec["interval"] != float64(2) || rec["endsAfter"] != float64(5) {
			t.Errorf("unexpected recurrence: %s", raw)
		}
		if wd, _ := rec["byWeekday"].([]any); len(wd) != 2 {
			t.Errorf("unexpected byWeekday: %s", raw)
		}
		if labels, _ := sent["waLabelIds"].([]any); len(labels) != 1 {
			t.Errorf("unexpected waLabelIds: %s", raw)
		}
		w.WriteHeader(http.StatusCreated)
		jsonResponse(w, `{"success":true,"data":{"id":"sm_2","waMessageKey":null,"to":"c1@c.us","recipientName":"Alice","type":"text","text":"hi","mediaUrl":null,"mediaKind":null,"pollQuestion":null,"pollOptions":null,"pollAllowMultiple":null,"status":"pending","sendAt":"2026-12-01T09:00:00Z","createdAt":"2026-01-01T00:00:00Z","confirmedAt":null,"receivedAt":null,"readAt":null,"playedAt":null,"failedAt":null,"failureReason":null,"secret":null,"linkPreview":null,"recurrence":{"frequency":"weekly","interval":2,"byWeekday":["MO","FR"],"endsAfter":5,"endsOn":null,"timezone":"UTC"},"waLabels":{"requested":1,"applied":false,"error":"wa_business_required"},"warning":{"code":"no_engine_connected","message":"No engine was connected when scheduled."}}}`)
	})
	two, five := 2, 5
	got, err := c.ScheduledMessages.Create(context.Background(), "c1@c.us", SendMessageParams{
		Type: MessageTypeText, Text: "hi", From: "+15550001234", WaLabelIDs: []string{"1"},
		Recurrence: &RecurrenceInput{
			Frequency: RecurrenceFrequencyWeekly, Interval: &two,
			ByWeekday: []string{RecurrenceWeekdayMO, RecurrenceWeekdayFR}, EndsAfter: &five,
		},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.RecipientName == nil || *got.RecipientName != "Alice" || got.Recurrence == nil || got.Recurrence.Frequency != RecurrenceFrequencyWeekly {
		t.Fatalf("unexpected recurrence/recipient: %+v", got)
	}
	if got.WaLabels == nil || got.WaLabels.Applied || got.WaLabels.Error == nil || got.Warning == nil || got.Warning.Code != "no_engine_connected" {
		t.Fatalf("unexpected waLabels/warning: %+v / %+v", got.WaLabels, got.Warning)
	}
}

func TestScheduledMessagesUpdateMediaFileID(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var sent map[string]any
		_ = json.Unmarshal(raw, &sent)
		if sent["mediaFileId"] != "f_9" || sent["mediaMimeType"] != "video/mp4" {
			t.Errorf("unexpected body: %s", raw)
		}
		jsonResponse(w, `{"success":true,"data":`+smData("sm_3", "media")+`}`)
	})
	fid, mime := "f_9", "video/mp4"
	if _, err := c.ScheduledMessages.Update(context.Background(), "sm_3", UpdateScheduledMessageParams{MediaFileID: &fid, MediaMimeType: &mime}); err != nil {
		t.Fatalf("Update: %v", err)
	}
}
