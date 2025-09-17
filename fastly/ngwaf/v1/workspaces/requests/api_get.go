package requests

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v12/fastly"
)

// GetInput specifies the information needed for the Get() function to
// perform the operation.
type GetInput struct {
	// RequestID is the request identifier (required).
	RequestID *string
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Get retrieves the specified request.
func Get(ctx context.Context, c *fastly.Client, i *GetInput) (*Request, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.RequestID == nil {
		return nil, fastly.ErrMissingRequestID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "requests", *i.RequestID)

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var request *Request
	if err := json.NewDecoder(resp.Body).Decode(&request); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return request, nil
}
