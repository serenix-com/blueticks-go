package blueticks

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestChatsList(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/chats" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("query") != "alice" {
			t.Errorf("missing query param")
		}
		jsonResponse(w, `{"data":[{"id":"c1@c.us","name":"Alice","is_group":false}],"has_more":false}`)
	})
	got, err := c.Chats.List(context.Background(), ListChatsParams{Query: "alice"})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got.Data) != 1 {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestChatsListMessagesEncodesMessageTypes(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/v1/chats/") || !strings.HasSuffix(r.URL.Path, "/messages") {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("message_types") != "image,document" {
			t.Errorf("expected comma-joined message_types, got %q", r.URL.Query().Get("message_types"))
		}
		if r.URL.Query().Get("mode") != "history" {
			t.Errorf("expected mode=history, got %q", r.URL.Query().Get("mode"))
		}
		jsonResponse(w, `{"data":[],"has_more":false}`)
	})
	_, err := c.Chats.ListMessages(context.Background(), "c1@c.us", ListChatMessagesParams{
		Mode: "history", MessageTypes: []string{"image", "document"},
	})
	if err != nil {
		t.Fatalf("ListMessages: %v", err)
	}
}

func TestChatsActions(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/v1/chats/c1@c.us/mark_read":
			jsonResponse(w, `{"ok":true}`)
		case r.URL.Path == "/v1/chats/c1@c.us/open":
			jsonResponse(w, `{"ok":true}`)
		case r.URL.Path == "/v1/chats/c1@c.us/messages/k1/reactions":
			jsonResponse(w, `{"ok":true}`)
		case r.URL.Path == "/v1/chats/c1@c.us/messages/load_older":
			jsonResponse(w, `{"loaded":10}`)
		case r.URL.Path == "/v1/chats/message_acks":
			jsonResponse(w, `{"k1":{"ack":3}}`)
		default:
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	if _, err := c.Chats.MarkRead(ctx, "c1@c.us"); err != nil {
		t.Fatalf("MarkRead: %v", err)
	}
	if _, err := c.Chats.Open(ctx, "c1@c.us"); err != nil {
		t.Fatalf("Open: %v", err)
	}
	if _, err := c.Chats.React(ctx, "c1@c.us", "k1", "🔥"); err != nil {
		t.Fatalf("React: %v", err)
	}
	if _, err := c.Chats.LoadOlderMessages(ctx, "c1@c.us"); err != nil {
		t.Fatalf("LoadOlderMessages: %v", err)
	}
	if _, err := c.Chats.BatchMessageAcks(ctx, []string{"k1"}); err != nil {
		t.Fatalf("BatchMessageAcks: %v", err)
	}
}

func TestChatsUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Chats.List(context.Background(), ListChatsParams{})
	expectAuthError(t, err)
}
