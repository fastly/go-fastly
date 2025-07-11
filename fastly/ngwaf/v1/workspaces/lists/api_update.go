package lists

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// UpdateInput specifies the information needed for the Update() function to
// perform the operation.
type UpdateInput struct {
	// Description is the description of the list.
	Description *string `json:"description,omitempty"`
	// Entries are the entries of the list.
	Entries *[]string `json:"entries,omitempty"`
	// ListID is the ID of the list to update (required).
	ListID *string
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Update updates the specified workspace.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*List, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.ListID == nil {
		return nil, fastly.ErrMissingListID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "lists", *i.ListID)

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var list *List
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}
	return list, nil
}
