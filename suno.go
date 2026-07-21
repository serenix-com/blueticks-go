package blueticks

import (
	"context"
	"net/http"
	"net/url"
)

// SunoResource exposes /v1/suno — AI song generation backed by the connected
// Suno account.
type SunoResource struct {
	client *Client
}

// Account returns the connected Suno account's credit balance and plan.
func (r *SunoResource) Account(ctx context.Context) (*SunoAccount, error) {
	body, err := r.client.do(ctx, requestOptions{method: http.MethodGet, path: "/v1/suno/account"})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data SunoAccount `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// GenerateSong submits a generation request and returns the batch job plus its
// (typically two) clips. Poll each clip with GetSong until it reaches a terminal
// status ("complete" or "error").
func (r *SunoResource) GenerateSong(ctx context.Context, params GenerateSongParams) (*GenerateSongResult, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/suno/songs",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data GenerateSongResult `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// GetSong returns the current status of a generated clip by id.
func (r *SunoResource) GetSong(ctx context.Context, id string) (*SongClip, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/suno/songs/" + url.PathEscape(id),
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data SongClip `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// CreateUpload uploads reference audio (by URL or base64) for cover / transform
// generations. Pass the returned UploadID as GenerateSongParams.UploadID.
func (r *SunoResource) CreateUpload(ctx context.Context, params CreateUploadParams) (*CreateUploadResult, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/suno/uploads",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data CreateUploadResult `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}
