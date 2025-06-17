package workspacethresholds

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// ListInput specifies the information needed for the List() function to perform
// the operation.
type ListInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// List retrieves a list of workspaces, with optional filtering and pagination.
func List(c *fastly.Client, i *ListInput) (*WorkspaceThresholds, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "thresholds")

	resp, err := c.Get(path, fastly.CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wts *WorkspaceThresholds
	if err := json.NewDecoder(resp.Body).Decode(&wts); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return wts, nil
}
