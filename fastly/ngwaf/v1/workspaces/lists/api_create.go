package lists

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// CreateInput specifies the information needed for the Create() function to
// perform the operation.
type CreateInput struct {
	// Description is the description of the list.
	Description *string `json:"description,omitempty"`
	// Entries are the entries of the list (required).
	Entries *[]string `json:"entries"`
	// Name is the name of the list (required).
	Name *string `json:"name"`
	// Type is the type of the list. Must be one of string | wildcard | ip | country | signal (required).
	Type *string `json:"type"`
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Create creates a new list.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*List, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.Entries == nil {
		return nil, fastly.ErrMissingEntries
	}
	if i.Name == nil {
		return nil, fastly.ErrMissingName
	}
	if i.Type == nil {
		return nil, fastly.ErrMissingType
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "lists")

	resp, err := c.PostJSON(ctx, path, i, fastly.CreateRequestOptions())
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
