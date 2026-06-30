package blueticks

import (
	"context"
	"encoding/json"
	"io"
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
		jsonResponse(w, `{"data":[{"id":"c1@c.us","name":"Alice","isGroup":false}],"hasMore":false}`)
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
		jsonResponse(w, `{"data":[{"chatId":"c1@c.us","isAdmin":true,"isSuperAdmin":false}],"hasMore":false,"nextCursor":null}`)
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
		if r.URL.Query().Get("messageTypes") != "image,document" {
			t.Errorf("expected comma-joined messageTypes, got %q", r.URL.Query().Get("messageTypes"))
		}
		if r.URL.Query().Get("order") != "asc" {
			t.Errorf("expected order=asc, got %q", r.URL.Query().Get("order"))
		}
		jsonResponse(w, `{"data":[],"hasMore":false}`)
	})
	_, err := c.Chats.ListMessages(context.Background(), "c1@c.us", ListChatMessagesParams{
		Order: "asc", MessageTypes: []string{"image", "document"},
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
			jsonResponse(w, `{"chatId":"c1@c.us"}`)
		case r.URL.Path == "/v1/messages/reactions/k1":
			jsonResponse(w, `{"ok":true}`)
		case r.URL.Path == "/v1/messages/ack/k1":
			if got := r.URL.Query().Get("chatId"); got != "c1@c.us" {
				t.Errorf("expected chatId query, got %q", got)
			}
			jsonResponse(w, `{"ack":3}`)
		case r.URL.Path == "/v1/chats/c1@c.us/messages/load_older":
			jsonResponse(w, `{"totalMessages":120,"added":40,"canLoadMore":true}`)
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
	rx, err := c.Chats.React(ctx, "k1", "", "🔥")
	if err != nil || !rx.Ok {
		t.Fatalf("React: %v %+v", err, rx)
	}
	ack, err := c.Chats.GetMessageAck(ctx, "k1", "c1@c.us")
	if err != nil || ack.Ack == nil || *ack.Ack != 3 {
		t.Fatalf("GetMessageAck: %v %+v", err, ack)
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
	if err := decodeJSON([]byte(`{}`), &BatchMessageAcksResponse{}); err != nil {
		t.Fatalf("BatchMessageAcksResponse: %v", err)
	}
	if err := decodeJSON([]byte(`{}`), &ParticipantsPage{}); err != nil {
		t.Fatalf("ParticipantsPage: %v", err)
	}
}

func TestChatsSendMessageText(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/chats/c1@c.us/messages" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Idempotency-Key"); got != "idem-text-1" {
			t.Errorf("expected Idempotency-Key header, got %q", got)
		}
		raw, _ := io.ReadAll(r.Body)
		var body map[string]interface{}
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["type"] != "text" || body["text"] != "hello" {
			t.Errorf("unexpected body: %+v", body)
		}
		if _, present := body["to"]; present {
			t.Errorf("body must not include `to` (derived from path)")
		}
		if _, present := body["sendAt"]; present {
			t.Errorf("body must not include `sendAt` (fire-and-forget)")
		}
		jsonResponse(w, `{"id":null,"key":"3EB0_text_1","to":"c1@c.us","type":"text","text":"hello","status":"sent","createdAt":"2026-05-28T00:00:00Z"}`)
	})
	got, err := c.Chats.SendMessage(context.Background(), "c1@c.us", SendMessageInChatParams{
		Type:           MessageTypeText,
		Text:           "hello",
		IdempotencyKey: "idem-text-1",
	})
	if err != nil {
		t.Fatalf("SendMessage: %v", err)
	}
	if got.Key == nil || *got.Key != "3EB0_text_1" {
		t.Fatalf("expected key=3EB0_text_1, got %+v", got.Key)
	}
	if got.Type != "text" {
		t.Fatalf("expected type=text, got %q", got.Type)
	}
}

func TestChatsSendMessageMedia(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/chats/c1@c.us/messages" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		var body map[string]interface{}
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["type"] != "media" {
			t.Errorf("expected type=media, got %v", body["type"])
		}
		media, ok := body["media"].(map[string]interface{})
		if !ok || media["url"] != "https://example.com/cat.jpg" || media["kind"] != "image" {
			t.Errorf("unexpected media block: %+v", body["media"])
		}
		jsonResponse(w, `{"id":null,"key":"3EB0_media_1","to":"c1@c.us","type":"media","mediaUrl":"https://example.com/cat.jpg","mediaKind":"image","status":"sent","createdAt":"2026-05-28T00:00:00Z"}`)
	})
	got, err := c.Chats.SendMessage(context.Background(), "c1@c.us", SendMessageInChatParams{
		Type: MessageTypeMedia,
		Media: &SendMediaObject{
			URL:     "https://example.com/cat.jpg",
			Kind:    "image",
			Caption: "neko",
		},
	})
	if err != nil {
		t.Fatalf("SendMessage: %v", err)
	}
	if got.MediaURL == nil || *got.MediaURL != "https://example.com/cat.jpg" {
		t.Fatalf("unexpected mediaUrl: %+v", got.MediaURL)
	}
}

func TestChatsSendMessagePoll(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/chats/c1@c.us/messages" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		var body map[string]interface{}
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		poll, ok := body["poll"].(map[string]interface{})
		if !ok || poll["question"] != "Pizza?" {
			t.Errorf("unexpected poll block: %+v", body["poll"])
		}
		opts, ok := poll["options"].([]interface{})
		if !ok || len(opts) != 2 || opts[0] != "yes" || opts[1] != "no" {
			t.Errorf("unexpected poll.options: %+v", poll["options"])
		}
		jsonResponse(w, `{"id":null,"key":"3EB0_poll_1","to":"c1@c.us","type":"poll","pollQuestion":"Pizza?","status":"sent","createdAt":"2026-05-28T00:00:00Z"}`)
	})
	got, err := c.Chats.SendMessage(context.Background(), "c1@c.us", SendMessageInChatParams{
		Type: MessageTypePoll,
		Poll: &SendPollObject{
			Question: "Pizza?",
			Options:  []string{"yes", "no"},
		},
	})
	if err != nil {
		t.Fatalf("SendMessage: %v", err)
	}
	if got.PollQuestion == nil || *got.PollQuestion != "Pizza?" {
		t.Fatalf("unexpected pollQuestion: %+v", got.PollQuestion)
	}
}

func TestChatsSendMessageUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Chats.SendMessage(context.Background(), "c1@c.us", SendMessageInChatParams{
		Type: MessageTypeText,
		Text: "hi",
	})
	expectAuthError(t, err)
}

func TestChatsUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Chats.List(context.Background(), ListChatsParams{})
	expectAuthError(t, err)
}
