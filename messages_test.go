package blueticks

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestMessagesList(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/messages" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("messageTypes"); got != "image,document" {
			t.Errorf("expected comma-joined messageTypes, got %q", got)
		}
		if got := r.URL.Query().Get("chatId"); got != "c1@c.us" {
			t.Errorf("expected chatId, got %q", got)
		}
		jsonResponse(w, `{"success":true,"data":[{"waMessageKey":{"fromMe":false,"id":"k1"},"chatId":"c1@c.us","from":"c1@c.us","type":"chat","fromMe":false,"text":"hi"}],"limit":50,"skip":0,"total":1,"hasMore":false}`)
	})
	got, err := c.Messages.List(context.Background(), ListMessagesParams{
		ChatID: "c1@c.us", MessageTypes: []string{"image", "document"},
	})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got.Data) != 1 || got.Data[0].ChatID != "c1@c.us" || got.Data[0].WaMessageKey.ID != "k1" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestMessagesGet(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/messages/k1" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("chatId") != "c1@c.us" {
			t.Errorf("expected chatId query")
		}
		jsonResponse(w, `{"success":true,"data":{"waMessageKey":{"fromMe":true,"id":"k1"},"chatId":"c1@c.us","from":"me","type":"chat","fromMe":true,"text":"yo"}}`)
	})
	got, err := c.Messages.Get(context.Background(), "k1", GetMessageParams{ChatID: "c1@c.us"})
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Text == nil || *got.Text != "yo" || !got.FromMe {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestMessagesSend(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/messages/c1@c.us" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Idempotency-Key"); got != "idem-1" {
			t.Errorf("expected Idempotency-Key header, got %q", got)
		}
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["type"] != "text" || body["text"] != "hello" {
			t.Errorf("unexpected body: %+v", body)
		}
		w.WriteHeader(http.StatusCreated)
		jsonResponse(w, `{"success":true,"data":{"id":"3EB0","waMessageKey":{"fromMe":true,"id":"3EB0"},"to":"c1@c.us","type":"text","text":"hello","mediaUrl":null,"mediaKind":null,"pollQuestion":null,"pollOptions":null,"pollAllowMultiple":null,"status":"sent","sendAt":null,"createdAt":"2026-05-28T00:00:00Z","confirmedAt":null,"receivedAt":null,"readAt":null,"playedAt":null,"failedAt":null,"failureReason":null,"secret":null}}`)
	})
	got, err := c.Messages.Send(context.Background(), "c1@c.us", SendInChatParams{
		Type: MessageTypeText, Text: "hello", IdempotencyKey: "idem-1",
	})
	if err != nil {
		t.Fatalf("Send: %v", err)
	}
	if got.ID != "3EB0" || got.Status != MessageStatusSent {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestMessagesAckAndBatch(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/messages/ack/k1":
			jsonResponse(w, `{"success":true,"data":{"ack":3}}`)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/messages/acks":
			jsonResponse(w, `{"success":true,"data":[{"key":"k1","ack":3,"found":true}],"limit":50,"skip":0,"total":1}`)
		default:
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	ack, err := c.Messages.GetAck(ctx, "k1", MessageAckParams{})
	if err != nil || ack.Ack == nil || *ack.Ack != 3 {
		t.Fatalf("GetAck: %v %+v", err, ack)
	}
	batch, err := c.Messages.BatchAcks(ctx, BatchMessageAcksParams{MessageKeys: []string{"k1"}})
	if err != nil || len(batch.Data) != 1 || !batch.Data[0].Found {
		t.Fatalf("BatchAcks: %v %+v", err, batch)
	}
}

func TestMessagesMediaAndLoadOlder(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/messages/media/k1":
			jsonResponse(w, `{"success":true,"data":{"url":"https://ex.com/m.jpg","mimetype":"image/jpeg"}}`)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/messages/load_older/c1@c.us":
			jsonResponse(w, `{"success":true,"data":{"totalMessages":120,"added":40,"canLoadMore":true,"historyUnavailable":false,"error":null}}`)
		default:
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	media, err := c.Messages.GetMedia(ctx, "k1", GetMediaParams{})
	if err != nil || media.URL == nil || *media.URL != "https://ex.com/m.jpg" {
		t.Fatalf("GetMedia: %v %+v", err, media)
	}
	lo, err := c.Messages.LoadOlder(ctx, "c1@c.us")
	if err != nil || !lo.CanLoadMore || lo.Added == nil || *lo.Added != 40 {
		t.Fatalf("LoadOlder: %v %+v", err, lo)
	}
}

func TestMessagesPinUnpinReactPinned(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/messages/pin/k1":
			jsonResponse(w, `{"success":true,"data":{"ok":true}}`)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/messages/unpin/k1":
			jsonResponse(w, `{"success":true,"data":{"ok":true}}`)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/messages/reactions/k1":
			jsonResponse(w, `{"success":true,"data":{"ok":true}}`)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/messages/pinned/c1@c.us":
			jsonResponse(w, `{"success":true,"data":[{"key":"k1","chatId":"c1@c.us","text":"pinned"}],"limit":50,"skip":0,"total":1}`)
		default:
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	dur := 86400
	if pin, err := c.Messages.Pin(ctx, "k1", PinMessageParams{Duration: &dur}); err != nil || !pin.Ok {
		t.Fatalf("Pin: %v %+v", err, pin)
	}
	if unpin, err := c.Messages.Unpin(ctx, "k1", UnpinMessageParams{}); err != nil || !unpin.Ok {
		t.Fatalf("Unpin: %v %+v", err, unpin)
	}
	if rx, err := c.Messages.React(ctx, "k1", ReactParams{Emoji: "🔥"}); err != nil || !rx.Ok {
		t.Fatalf("React: %v %+v", err, rx)
	}
	pinned, err := c.Messages.ListPinned(ctx, "c1@c.us")
	if err != nil || len(pinned.Data) != 1 || pinned.Data[0].Key != "k1" {
		t.Fatalf("ListPinned: %v %+v", err, pinned)
	}
}

func TestMessagesDecodesEmptyTypedResponses(t *testing.T) {
	if err := decodeJSON([]byte(`{}`), &Message{}); err != nil {
		t.Fatalf("Message: %v", err)
	}
	if err := decodeJSON([]byte(`{}`), &Media{}); err != nil {
		t.Fatalf("Media: %v", err)
	}
	if err := decodeJSON([]byte(`{}`), &LoadOlderResult{}); err != nil {
		t.Fatalf("LoadOlderResult: %v", err)
	}
	if err := decodeJSON([]byte(`{}`), &MessageAck{}); err != nil {
		t.Fatalf("MessageAck: %v", err)
	}
	if err := decodeJSON([]byte(`{}`), &PinnedMessagesPage{}); err != nil {
		t.Fatalf("PinnedMessagesPage: %v", err)
	}
}

func TestMessagesUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Messages.List(context.Background(), ListMessagesParams{})
	expectAuthError(t, err)
}

func TestMessagesListNewFilters(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q["queryAny"]; len(got) != 2 || got[0] != "invoice, final" || got[1] != "receipt" {
			t.Errorf("expected repeated queryAny, got %v", got)
		}
		if q.Get("mime") != "application/pdf" || q.Get("filenameContains") != "inv" || q.Get("sender") != "+1555" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		if q.Get("hasMedia") != "false" {
			t.Errorf("expected hasMedia=false to be sent, got %q", q.Get("hasMedia"))
		}
		jsonResponse(w, `{"success":true,"data":[{"waMessageKey":{"fromMe":true,"id":"k1"},"chatId":"c1@c.us","from":"me","type":"chat","fromMe":true,"source":"backup","secret":"tag-1"}],"limit":50,"skip":0,"total":1,"hasMore":false,"coverage":{"oldestLoadedTimestamp":1700000000000,"oldestLoadedIso":"2023-11-14T22:13:20.000Z"}}`)
	})
	no := false
	got, err := c.Messages.List(context.Background(), ListMessagesParams{
		QueryAny: []string{"invoice, final", "receipt"}, Mime: "application/pdf",
		FilenameContains: "inv", HasMedia: &no, Sender: "+1555",
	})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if got.Coverage == nil || got.Coverage.OldestLoadedIso == nil || *got.Coverage.OldestLoadedIso != "2023-11-14T22:13:20.000Z" {
		t.Fatalf("coverage: %+v", got.Coverage)
	}
	if got.Coverage.OldestLoadedTimestamp == nil || *got.Coverage.OldestLoadedTimestamp != 1700000000000 {
		t.Fatalf("coverage ts: %+v", got.Coverage)
	}
	m := got.Data[0]
	if m.Source == nil || *m.Source != "backup" || m.Secret == nil || *m.Secret != "tag-1" {
		t.Fatalf("source/secret: %+v", m)
	}
}

func TestMessagesLoadOlderWithParams(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/messages/load_older/c1@c.us" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Fatalf("decode body: %v (%s)", err, raw)
		}
		if body["pages"] != float64(3) || body["until_date"] != "2025-01-01T00:00:00Z" {
			t.Errorf("unexpected body: %s", raw)
		}
		jsonResponse(w, `{"success":true,"data":{"totalMessages":120,"added":40,"canLoadMore":true,"historyUnavailable":false,"error":null,"pagesExecuted":3,"coverage":{"oldestLoadedTimestamp":1700000000000,"oldestLoadedIso":"2023-11-14T22:13:20.000Z"}}}`)
	})
	pages := 3
	got, err := c.Messages.LoadOlderWithParams(context.Background(), "c1@c.us", LoadOlderParams{Pages: &pages, UntilDate: "2025-01-01T00:00:00Z"})
	if err != nil {
		t.Fatalf("LoadOlderWithParams: %v", err)
	}
	if got.PagesExecuted != 3 || got.Coverage.OldestLoadedIso == nil || !got.CanLoadMore {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestMessagesLoadOlderNumericUntilDate(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		_ = json.Unmarshal(raw, &body)
		if body["until_date"] != float64(1700000000) {
			t.Errorf("expected numeric until_date, got %s", raw)
		}
		if _, ok := body["pages"]; ok {
			t.Errorf("pages should be omitted: %s", raw)
		}
		jsonResponse(w, `{"success":true,"data":{"totalMessages":null,"added":0,"canLoadMore":false,"historyUnavailable":false,"error":null,"pagesExecuted":1,"coverage":{}}}`)
	})
	if _, err := c.Messages.LoadOlderWithParams(context.Background(), "c1@c.us", LoadOlderParams{UntilDate: 1700000000}); err != nil {
		t.Fatalf("LoadOlderWithParams: %v", err)
	}
}

func TestMessagesLoadOlderNoParamsSendsNoBody(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		if len(raw) != 0 {
			t.Errorf("expected empty body for plain LoadOlder, got %s", raw)
		}
		jsonResponse(w, `{"success":true,"data":{"totalMessages":1,"added":1,"canLoadMore":true,"historyUnavailable":false,"error":null,"pagesExecuted":1,"coverage":{}}}`)
	})
	got, err := c.Messages.LoadOlder(context.Background(), "c1@c.us")
	if err != nil {
		t.Fatalf("LoadOlder: %v", err)
	}
	if got.PagesExecuted != 1 || got.Coverage.OldestLoadedIso != nil {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestMessagesSendMentions(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var body struct {
			Mentions []string `json:"mentions"`
		}
		_ = json.Unmarshal(raw, &body)
		if len(body.Mentions) != 1 || body.Mentions[0] != "12025550101@c.us" {
			t.Errorf("unexpected body: %s", raw)
		}
		jsonResponse(w, `{"success":true,"data":{"id":"m1","to":"g1@g.us","type":"text","status":"sent","createdAt":"2026-01-01T00:00:00Z"}}`)
	})
	if _, err := c.Messages.Send(context.Background(), "g1@g.us", SendInChatParams{Type: MessageTypeText, Text: "hi @Jon", Mentions: []string{"12025550101@c.us"}}); err != nil {
		t.Fatalf("Send: %v", err)
	}
}
