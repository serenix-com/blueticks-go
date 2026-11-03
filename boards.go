package blueticks

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// BoardsResource exposes /v1/boards — agent-memory boards and their items.
type BoardsResource struct {
	client *Client
}

// List boards, newest first by default. Requires `boards:read`. Offset-paginated.
func (r *BoardsResource) List(ctx context.Context, params ListBoardsParams) (*BoardsPage, error) {
	q := url.Values{}
	if params.Order != "" {
		q.Set("order", params.Order)
	}
	if params.Skip > 0 {
		q.Set("skip", strconv.Itoa(params.Skip))
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/boards",
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	out := &BoardsPage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Create finds or creates the agent memory board for (AgentID, Domain).
//
// Pass Domain for a built-in board (tasks/leads/invoices) or EntityKind to open
// a NEW dynamic (custom:<slug>) board — Description is then required. Returns
// the existing board if one exists. Requires `boards:write`.
func (r *BoardsResource) Create(ctx context.Context, params CreateBoardParams) (*BoardDetail, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/boards",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data BoardDetail `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// Get retrieves a board with its columns, options and item count. Requires
// `boards:read`.
func (r *BoardsResource) Get(ctx context.Context, id string) (*BoardDetail, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/boards/" + url.PathEscape(id),
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data BoardDetail `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// ListItems lists a board's items (tasks). Set ActiveOnly to return only items
// updated within the retention window whose status is non-terminal. Requires
// `boards:read`. Offset-paginated.
func (r *BoardsResource) ListItems(ctx context.Context, boardID string, params ListBoardItemsParams) (*BoardItemsPage, error) {
	q := url.Values{}
	if params.Status != "" {
		q.Set("status", params.Status)
	}
	if params.Order != "" {
		q.Set("order", params.Order)
	}
	if params.ActiveOnly != nil {
		q.Set("active_only", strconv.FormatBool(*params.ActiveOnly))
	}
	if params.UpdatedSince != "" {
		q.Set("updated_since", params.UpdatedSince)
	}
	if params.Skip > 0 {
		q.Set("skip", strconv.Itoa(params.Skip))
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/boards/" + url.PathEscape(boardID) + "/items",
		query:  q,
	})
	if err != nil {
		return nil, err
	}
	out := &BoardItemsPage{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateItem creates one board item. Status and Priority are given as labels and
// resolved to the board's options server-side (unknown label -> 422). Requires
// `boards:write`.
func (r *BoardsResource) CreateItem(ctx context.Context, boardID string, params CreateBoardItemParams) (*BoardItem, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/boards/" + url.PathEscape(boardID) + "/items",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	return decodeBoardItem(body)
}

// BulkUpsertItems creates or merges many items in one call (1-500). Each item
// is fingerprinted and matched against the board's existing items: a match is
// patched, a miss is created. Requires `boards:write`.
func (r *BoardsResource) BulkUpsertItems(ctx context.Context, boardID string, params BulkBoardItemsParams) (*BulkBoardItemsResult, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/boards/" + url.PathEscape(boardID) + "/items/bulk",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data BulkBoardItemsResult `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// FindSimilarItems matches each candidate text against the board's items by
// normalized fingerprint + fuzzy token overlap. Use before creating to avoid
// duplicates. Requires `boards:read`.
func (r *BoardsResource) FindSimilarItems(ctx context.Context, boardID string, params FindSimilarBoardItemsParams) (*FindSimilarBoardItemsResult, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/boards/" + url.PathEscape(boardID) + "/items/find-similar",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data FindSimilarBoardItemsResult `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// UpdateItem updates cells, soft-deletes (Deleted=true) or restores
// (Deleted=false) an item. Note is recorded in the item's history. Requires
// `boards:write`.
func (r *BoardsResource) UpdateItem(ctx context.Context, boardID, itemID string, params UpdateBoardItemParams) (*BoardItem, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPatch,
		path:   "/v1/boards/" + url.PathEscape(boardID) + "/items/" + url.PathEscape(itemID),
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	return decodeBoardItem(body)
}

func decodeBoardItem(body []byte) (*BoardItem, error) {
	var env struct {
		Data BoardItem `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}
