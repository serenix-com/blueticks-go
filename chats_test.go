package blueticks

import (
	"context"
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
