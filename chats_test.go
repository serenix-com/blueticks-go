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

func TestChatsListParticipants(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/chats/g1@g.us/participants" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"data":[{"chat_id":"c1@c.us","is_admin":true,"is_super_admin":false}],"has_more":false,"next_cursor":null}`)
	})
	got, err := c.Chats.ListParticipants(context.Background(), "g1@g.us", ListParticipantsParams{Limit: 50})
	if err != nil {
		t.Fatalf("ListParticipants: %v", err)
	}
	if len(got.Data) != 1 || got.Data[0].ChatID != "c1@c.us" || !got.Data[0].IsAdmin {
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
			jsonResponse(w, `{"chat_id":"c1@c.us"}`)
		case r.URL.Path == "/v1/chats/c1@c.us/messages/k1/reactions":
			jsonResponse(w, `{"ok":true}`)
		case r.URL.Path == "/v1/chats/c1@c.us/messages/k1/ack":
			jsonResponse(w, `{"ack":3}`)
		case r.URL.Path == "/v1/chats/c1@c.us/messages/k1/media_url":
			jsonResponse(w, `{"url":"https://example.com/m.jpg"}`)
		case r.URL.Path == "/v1/chats/c1@c.us/messages/load_older":
			jsonResponse(w, `{"total_messages":120,"added":40,"can_load_more":true}`)
		case r.URL.Path == "/v1/chats/message_acks":
			jsonResponse(w, `{"data":[{"key":"k1","ack":3}]}`)
		default:
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()

	mr, err := c.Chats.MarkRead(ctx, "c1@c.us")
	if err != nil || !mr.Ok {
		t.Fatalf("MarkRead: %v %+v", err, mr)
	}
	op, err := c.Chats.Open(ctx, "c1@c.us")
	if err != nil || op.ChatID != "c1@c.us" {
		t.Fatalf("Open: %v %+v", err, op)
	}
	rx, err := c.Chats.React(ctx, "c1@c.us", "k1", "🔥")
	if err != nil || !rx.Ok {
		t.Fatalf("React: %v %+v", err, rx)
	}
	ack, err := c.Chats.GetMessageAck(ctx, "c1@c.us", "k1")
	if err != nil || ack.Ack == nil || *ack.Ack != 3 {
		t.Fatalf("GetMessageAck: %v %+v", err, ack)
	}
	mu, err := c.Chats.GetMediaURL(ctx, "c1@c.us", "k1")
	if err != nil || mu.URL == nil || *mu.URL != "https://example.com/m.jpg" {
		t.Fatalf("GetMediaURL: %v %+v", err, mu)
	}
	lo, err := c.Chats.LoadOlderMessages(ctx, "c1@c.us")
	if err != nil || !lo.CanLoadMore || lo.Added == nil || *lo.Added != 40 {
		t.Fatalf("LoadOlderMessages: %v %+v", err, lo)
	}
	bm, err := c.Chats.BatchMessageAcks(ctx, BatchMessageAcksParams{MessageKeys: []string{"k1"}})
	if err != nil || len(bm.Data) != 1 {
		t.Fatalf("BatchMessageAcks: %v %+v", err, bm)
	}
	if bm.Data[0].Key != "k1" || bm.Data[0].Ack == nil || *bm.Data[0].Ack != 3 {
		t.Fatalf("BatchMessageAcks entry: %+v", bm.Data[0])
	}
}

// TestChatsDecodesEmptyTypedResponses locks in encoding/json's lenient
// behavior for the new typed envelopes — an empty `{}` decodes to the
// zero value without panicking.
func TestChatsDecodesEmptyTypedResponses(t *testing.T) {
	if err := decodeJSON([]byte(`{}`), &OkResponse{}); err != nil {
		t.Fatalf("OkResponse: %v", err)
	}
	if err := decodeJSON([]byte(`{}`), &ChatRef{}); err != nil {
		t.Fatalf("ChatRef: %v", err)
	}
	if err := decodeJSON([]byte(`{}`), &MessageAck{}); err != nil {
		t.Fatalf("MessageAck: %v", err)
	}
	if err := decodeJSON([]byte(`{}`), &LoadOlderMessagesResponse{}); err != nil {
		t.Fatalf("LoadOlderMessagesResponse: %v", err)
	}
	if err := decodeJSON([]byte(`{}`), &MediaURLResponse{}); err != nil {
		t.Fatalf("MediaURLResponse: %v", err)
	}
	if err := decodeJSON([]byte(`{}`), &BatchMessageAcksResponse{}); err != nil {
		t.Fatalf("BatchMessageAcksResponse: %v", err)
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
