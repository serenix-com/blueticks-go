package blueticks

import (
	"context"
	"net/http"
)

// UploadsResource exposes /v1/uploads: signed links for sending files.
type UploadsResource struct {
	client *Client
}

// Create requests a signed upload link for one file. PUT the raw bytes to
// FileUpload.UploadURL with exactly FileUpload.Headers, then send with type
// "media" and mediaFileId set to FileUpload.FileID. The link expires after about
// 15 minutes; an unsent file is deleted about 24 hours after creation.
func (r *UploadsResource) Create(ctx context.Context, params CreateFileUploadParams) (*FileUpload, error) {
	body, err := r.client.do(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/v1/uploads",
		body:   params,
	})
	if err != nil {
		return nil, err
	}
	var env struct {
		Data FileUpload `json:"data"`
	}
	if err := decodeJSON(body, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}
