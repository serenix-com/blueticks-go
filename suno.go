package blueticks

import (
	"context"
	"net/http"
	"net/url"
)

// SunoResource exposes /v1/suno/* — Suno-backed song generation, reference
// uploads, and account status.
type SunoResource struct {
	client *Client
}

// Generate submits a song generation to Suno.
//
// Submit a song generation to Suno from `lyrics` + `style`, optionally
// steering `vocalGender`, `weirdness`, `styleInfluence`, and a reference
// recording (`uploadId` from SunoResource.Upload, with `audioInfluence`
// controlling how closely the cover follows it). Returns two clip variants —
// poll each with GetSong until Status is "complete". Requires `suno:write`.
func (r *SunoResource) Generate(ctx context.Context, params SunoGenerateParams) (*SunoGenerateResult, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/suno/songs",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	out := &SunoGenerateResult{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetSong polls a single generated clip by id.
//
// When Status is "complete", AudioURL (MP3) and ImageURL are populated.
// Requires `suno:read`.
func (r *SunoResource) GetSong(ctx context.Context, id string) (*SunoClip, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/suno/songs/" + url.PathEscape(id),
	})
	if err != nil {
		return nil, err
	}
	out := &SunoClip{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Upload uploads reference audio to Suno.
//
// Upload an audio recording (via AudioURL or base64 AudioBase64). Returns an
// UploadID to pass to Generate as SunoGenerateParams.UploadID to
// cover/transform it. Max 500 MB. Requires `suno:write`.
func (r *SunoResource) Upload(ctx context.Context, params SunoUploadParams) (*SunoUpload, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/suno/uploads",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	out := &SunoUpload{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetAccount returns remaining credits, monthly usage, and plan on the
// connected Suno account.
//
// Requires `suno:read`.
func (r *SunoResource) GetAccount(ctx context.Context) (*SunoAccount, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/v1/suno/account",
	})
	if err != nil {
		return nil, err
	}
	out := &SunoAccount{}
	if err := decodeJSON(body, out); err != nil {
		return nil, err
	}
	return out, nil
}
