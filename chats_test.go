package blueticks

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestChatsList(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/chats" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("searchToken") != "alice" {
			t.Errorf("missing searchToken param")
		}
		jsonResponse(w, `{"success":true,"data":[{"chatId":"c1@c.us","name":"Alice","chatType":"contact","markedUnread":false,"lastMessageAt":null,"unreadCount":0}],"limit":50,"skip":0,"total":1}`)
	})
	got, err := c.Chats.List(context.Background(), ListChatsParams{SearchToken: "alice"})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got.Data) != 1 || got.Data[0].ChatType != ChatTypeContact {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestChatsGet(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/chats/c1@c.us" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"success":true,"data":{"chatId":"c1@c.us","name":"Alice","chatType":"contact","markedUnread":false,"lastMessageAt":null,"unreadCount":2}}`)
	})
	got, err := c.Chats.Get(context.Background(), "c1@c.us")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.ID != "c1@c.us" || got.UnreadCount == nil || *got.UnreadCount != 2 {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestChatsListParticipants(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/chats/g1@g.us/participants" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"success":true,"data":[{"chatId":"c1@c.us","name":"Alice","isAdmin":true,"isSuperAdmin":false}],"limit":100,"skip":0,"total":1}`)
	})
	got, err := c.Chats.ListParticipants(context.Background(), "g1@g.us", ListParticipantsParams{Limit: 50})
	if err != nil {
		t.Fatalf("ListParticipants: %v", err)
	}
	if len(got.Data) != 1 || got.Data[0].ChatID != "c1@c.us" || !got.Data[0].IsAdmin {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestChatsActions(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/chats/c1@c.us/archive", "/v1/chats/c1@c.us/unarchive", "/v1/chats/c1@c.us/mark_read":
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			jsonResponse(w, `{"success":true,"data":{"ok":true}}`)
		default:
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	if ar, err := c.Chats.Archive(ctx, "c1@c.us"); err != nil || !ar.Ok {
		t.Fatalf("Archive: %v %+v", err, ar)
	}
	if ua, err := c.Chats.Unarchive(ctx, "c1@c.us"); err != nil || !ua.Ok {
		t.Fatalf("Unarchive: %v %+v", err, ua)
	}
	if mr, err := c.Chats.MarkRead(ctx, "c1@c.us"); err != nil || !mr.Ok {
		t.Fatalf("MarkRead: %v %+v", err, mr)
	}
}

func TestChatsDecodesEmptyTypedResponses(t *testing.T) {
	if err := decodeJSON([]byte(`{}`), &Chat{}); err != nil {
		t.Fatalf("Chat: %v", err)
	}
	if err := decodeJSON([]byte(`{}`), &OkResult{}); err != nil {
		t.Fatalf("OkResult: %v", err)
	}
	if err := decodeJSON([]byte(`{}`), &ParticipantsPage{}); err != nil {
		t.Fatalf("ParticipantsPage: %v", err)
	}
}

func TestChatsUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Chats.List(context.Background(), ListChatsParams{})
	expectAuthError(t, err)
}

func TestChatsListKindsAndSince(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q["kinds"]; len(got) != 2 || got[0] != "group" || got[1] != "newsletter" {
			t.Errorf("expected repeated kinds, got %v", got)
		}
		if q.Get("since") != "2026-01-01T00:00:00+02:00" {
			t.Errorf("unexpected since: %q", q.Get("since"))
		}
		jsonResponse(w, `{"success":true,"data":[],"limit":50,"skip":0,"total":0}`)
	})
	if _, err := c.Chats.List(context.Background(), ListChatsParams{Kinds: []string{"group", "newsletter"}, Since: "2026-01-01T00:00:00+02:00"}); err != nil {
		t.Fatalf("List: %v", err)
	}
}

func TestChatsListWaLabels(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/chats/c1@c.us/wa-labels" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"success":true,"data":[{"id":"1","name":"VIP","color":"#00ff00"}]}`)
	})
	got, err := c.Chats.ListWaLabels(context.Background(), "c1@c.us")
	if err != nil {
		t.Fatalf("ListWaLabels: %v", err)
	}
	if len(got) != 1 || got[0].ID != "1" || got[0].Name != "VIP" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestChatsUpdateWaLabels(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/v1/chats/c1@c.us/wa-labels" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		var sent map[string]any
		_ = json.Unmarshal(raw, &sent)
		ids, ok := sent["labelIds"].([]any)
		if !ok || len(ids) != 0 {
			t.Errorf("an empty labelIds must go out as [] (clears labels), got %s", raw)
		}
		jsonResponse(w, `{"success":true,"data":{"ok":true}}`)
	})
	got, err := c.Chats.UpdateWaLabels(context.Background(), "c1@c.us", UpdateChatWaLabelsParams{LabelIDs: []string{}})
	if err != nil {
		t.Fatalf("UpdateWaLabels: %v", err)
	}
	if !got.Ok {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestChatsListEngineTargeting(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/chats" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		q := r.URL.Query()
		want := map[string]string{"engine": "gateway", "engineId": "eng_1", "numberId": "num_1", "phone": "+15550100", "fanoutLimit": "3"}
		for k, v := range want {
			if q.Get(k) != v {
				t.Errorf("query %s=%q, want %q", k, q.Get(k), v)
			}
		}
		jsonResponse(w, `{"success":true,"data":[],"limit":50,"skip":0,"total":0,"answeredBy":{"engineId":"eng_1","numberId":"num_1","phoneNumber":"+15550100"},"nextCursor":"cur_2"}`)
	})
	got, err := c.Chats.List(context.Background(), ListChatsParams{EngineTarget: EngineTarget{
		Engine: EngineTypeGateway, EngineID: "eng_1", NumberID: "num_1", Phone: "+15550100", FanoutLimit: 3,
	}})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if got.AnsweredBy == nil || got.AnsweredBy.EngineID != "eng_1" || got.NextCursor == nil || *got.NextCursor != "cur_2" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestChatsParticipantsEngineTargeting(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/chats/c1@g.us/participants" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("engine") != "regular" || r.URL.Query().Get("fanoutLimit") != "2" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		jsonResponse(w, `{"success":true,"data":[{"chatId":"p1@c.us","isAdmin":false}],"limit":50,"skip":0,"total":1}`)
	})
	got, err := c.Chats.ListParticipants(context.Background(), "c1@g.us", ListParticipantsParams{
		EngineTarget: EngineTarget{Engine: EngineTypeRegular, FanoutLimit: 2},
	})
	if err != nil {
		t.Fatalf("ListParticipants: %v", err)
	}
	if len(got.Data) != 1 || got.Data[0].ChatID != "p1@c.us" {
		t.Fatalf("unexpected: %+v", got)
	}
}
