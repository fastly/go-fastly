package workspaces

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// UpdateInput specifies the information needed to update a workspace.
type UpdateInput struct {
	// WorkspaceID is the ID of the bot management workspace (required).
	WorkspaceID *string `json:"-"`
	// Description is a freeform descriptive note.
	Description *string `json:"description,omitempty"`
	// Name is the display name of the workspace.
	Name *string `json:"name,omitempty"`
	// ProtectionMode is the protection mode of the workspace. Must be one of
	// `off`, `log`, or `block`.
	ProtectionMode *string `json:"protection_mode,omitempty"`
}

// Update updates a specified workspace.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*Workspace, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	path := fastly.ToSafeURL("bot-management", "v1", "workspaces", *i.WorkspaceID)

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
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
