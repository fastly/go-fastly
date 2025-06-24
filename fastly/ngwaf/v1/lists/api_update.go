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
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// Description is the description of the list.
	Description *string `json:"description"`
	// Entries are the entries of the list.
	Entries *[]string `json:"entries"`
	// ListID is the ID of the list to update (required).
	ListID *string
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Update updates the specified workspace.
func Update(c *fastly.Client, i *UpdateInput) (*List, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.ListID == nil {
		return nil, fastly.ErrMissingListID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "lists", *i.ListID)

	resp, err := c.PatchJSON(path, i, fastly.CreateRequestOptions(i.Context))
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
