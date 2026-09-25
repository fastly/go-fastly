package workspaces

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// GetInput specifies the information needed to get a workspace.
type GetInput struct {
	// WorkspaceID is the ID of the bot management workspace (required).
	WorkspaceID *string `json:"-"`
}

// Get retrieves a specified workspace.
func Get(ctx context.Context, c *fastly.Client, i *GetInput) (*Workspace, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	path := fastly.ToSafeURL("bot-management", "v1", "workspaces", *i.WorkspaceID)

	resp, err := c.GetJSON(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var workspace *Workspace
	if err := json.NewDecoder(resp.Body).Decode(&workspace); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return workspace, nil
}
