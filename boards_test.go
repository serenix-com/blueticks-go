package blueticks

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

const boardJSON = `{"board_id":"b1","name":"Tasks","domain":"tasks","agent_managed":true,"memory_retention_days":30,"item_count":2,"web_url":"https://app/b1","created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-02T00:00:00Z"`

const itemJSON = `{"item_id":"i1","text":"call Dana","status":{"label":"Open","bgcolor":"#fff"},"due_date":"2026-02-01","cells":{"k":1},"source":{"chat_id":"c1@c.us","extracted_at":123},"confidence":0.9,"last_seen_at":456,"web_url":"https://app/i1","created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-02T00:00:00Z"}`

func TestBoardsList(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/boards" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("order") != "asc" || r.URL.Query().Get("limit") != "5" || r.URL.Query().Get("skip") != "10" {
			t.Errorf("unexpected query: %v", r.URL.RawQuery)
		}
		jsonResponse(w, `{"success":true,"data":[`+boardJSON+`}],"limit":5,"skip":10,"total":11}`)
	})
	got, err := c.Boards.List(context.Background(), ListBoardsParams{Order: "asc", Limit: 5, Skip: 10})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got.Data) != 1 || got.Data[0].BoardID != "b1" || !got.Data[0].AgentManaged || got.Total != 11 {
		t.Fatalf("unexpected: %+v", got)
	}
	if got.Data[0].ItemCount == nil || *got.Data[0].ItemCount != 2 {
		t.Fatalf("item count: %+v", got.Data[0])
	}
}

func TestBoardsListUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Boards.List(context.Background(), ListBoardsParams{})
	expectAuthError(t, err)
}

func TestBoardsCreate(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/boards" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		_ = json.Unmarshal(raw, &body)
		if body["entity_kind"] != "orders" || body["description"] != "order tracking" {
			t.Errorf("unexpected body: %s", raw)
		}
		if _, ok := body["domain"]; ok {
			t.Errorf("domain should be omitted: %s", raw)
		}
		jsonResponse(w, `{"success":true,"data":`+boardJSON+`,"columns":[{"column_id":"col1","title":"Status","type":"status","options":[{"label":"Open","bgcolor":"#fff"}]}]}}`)
	})
	got, err := c.Boards.Create(context.Background(), CreateBoardParams{EntityKind: "orders", Description: "order tracking"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.BoardID != "b1" || len(got.Columns) != 1 || got.Columns[0].Options[0].Label != "Open" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestBoardsGet(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/boards/b1" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"success":true,"data":`+boardJSON+`,"columns":[]}}`)
	})
	got, err := c.Boards.Get(context.Background(), "b1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != "Tasks" || got.Domain == nil || *got.Domain != "tasks" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestBoardsListItems(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/boards/b1/items" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("active_only") != "true" || q.Get("status") != "Open" || q.Get("updated_since") != "2026-01-01" {
			t.Errorf("unexpected query: %v", r.URL.RawQuery)
		}
		jsonResponse(w, `{"success":true,"data":[`+itemJSON+`],"limit":50,"skip":0,"total":1}`)
	})
	active := true
	got, err := c.Boards.ListItems(context.Background(), "b1", ListBoardItemsParams{Status: "Open", ActiveOnly: &active, UpdatedSince: "2026-01-01"})
	if err != nil {
		t.Fatalf("ListItems: %v", err)
	}
	if len(got.Data) != 1 || got.Data[0].ItemID != "i1" || got.Data[0].Status == nil || got.Data[0].Status.Label != "Open" {
		t.Fatalf("unexpected: %+v", got)
	}
	if got.Data[0].Source == nil || got.Data[0].Source.ExtractedAt == nil || *got.Data[0].Source.ExtractedAt != 123 {
		t.Fatalf("source: %+v", got.Data[0].Source)
	}
}

func TestBoardsCreateItem(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/boards/b1/items" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		_ = json.Unmarshal(raw, &body)
		if body["text"] != "call Dana" || body["status"] != "Open" || body["allow_duplicate"] != true {
			t.Errorf("unexpected body: %s", raw)
		}
		w.WriteHeader(http.StatusCreated)
		jsonResponse(w, `{"success":true,"data":`+itemJSON+`}`)
	})
	dup := true
	got, err := c.Boards.CreateItem(context.Background(), "b1", CreateBoardItemParams{Text: "call Dana", Status: "Open", AllowDuplicate: &dup})
	if err != nil {
		t.Fatalf("CreateItem: %v", err)
	}
	if got.ItemID != "i1" || got.Confidence == nil || *got.Confidence != 0.9 {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestBoardsBulkUpsertItems(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/boards/b1/items/bulk" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		var body struct {
			Items []map[string]any `json:"items"`
			Mode  string           `json:"mode"`
		}
		_ = json.Unmarshal(raw, &body)
		if len(body.Items) != 2 || body.Mode != "upsert" {
			t.Errorf("unexpected body: %s", raw)
		}
		jsonResponse(w, `{"success":true,"data":{"created":[`+itemJSON+`],"updated":[],"unchanged":1,"board":{"board_id":"b1","name":"Tasks","web_url":"https://app/b1"}}}`)
	})
	got, err := c.Boards.BulkUpsertItems(context.Background(), "b1", BulkBoardItemsParams{
		Items: []CreateBoardItemParams{{Text: "a"}, {Text: "b"}}, Mode: "upsert",
	})
	if err != nil {
		t.Fatalf("BulkUpsertItems: %v", err)
	}
	if len(got.Created) != 1 || got.Unchanged != 1 || got.Board.BoardID != "b1" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestBoardsFindSimilarItems(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/boards/b1/items/find-similar" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"success":true,"data":{"matches":[{"candidate_index":0,"item_id":"i1","text":"call Dana","score":0.8}],"misses":[1]}}`)
	})
	got, err := c.Boards.FindSimilarItems(context.Background(), "b1", FindSimilarBoardItemsParams{
		Candidates: []FindSimilarCandidate{{Text: "call dana"}, {Text: "other"}},
	})
	if err != nil {
		t.Fatalf("FindSimilarItems: %v", err)
	}
	if len(got.Matches) != 1 || got.Matches[0].Score != 0.8 || len(got.Misses) != 1 || got.Misses[0] != 1 {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestBoardsUpdateItem(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/v1/boards/b1/items/i1" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		_ = json.Unmarshal(raw, &body)
		if body["deleted"] != true || body["note"] != "dup" {
			t.Errorf("unexpected body: %s", raw)
		}
		if _, ok := body["text"]; ok {
			t.Errorf("text should be omitted: %s", raw)
		}
		jsonResponse(w, `{"success":true,"data":`+itemJSON+`}`)
	})
	del, note := true, "dup"
	got, err := c.Boards.UpdateItem(context.Background(), "b1", "i1", UpdateBoardItemParams{Deleted: &del, Note: &note})
	if err != nil {
		t.Fatalf("UpdateItem: %v", err)
	}
	if got.ItemID != "i1" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestBoardsItemOpsUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Boards.UpdateItem(context.Background(), "b1", "i1", UpdateBoardItemParams{})
	expectAuthError(t, err)
	_, err = c.Boards.Get(context.Background(), "b1")
	expectAuthError(t, err)
}

func TestBoardsDecodeEmptyResponse(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) { jsonResponse(w, `{}`) })
	got, err := c.Boards.Get(context.Background(), "b1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.BoardID != "" || got.Columns != nil {
		t.Fatalf("expected zero value, got %+v", got)
	}
	page, err := c.Boards.List(context.Background(), ListBoardsParams{})
	if err != nil || page.Total != 0 {
		t.Fatalf("List empty: %v %+v", err, page)
	}
}
