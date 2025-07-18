package lists

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v11/fastly"
	"github.com/fastly/go-fastly/v11/fastly/ngwaf/v1/common"
)

// UpdateInput specifies the information needed for the Update()
// function to perform the operation.
type UpdateInput struct {
	// Description is the description of the list.
	Description *string `json:"description,omitempty"`
	// Entries are the entries of the list.
	Entries *[]string `json:"entries,omitempty"`
	// ListID is the ID of the list to update (required).
	ListID *string
	// Scope defines where the list is located, including its type (e.g.,
	// "workspace" or "account") and the specific IDs it applies to (required).
	Scope *common.Scope
}

// Update updates the specified list.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*List, error) {
	if i.ListID == nil {
		return nil, fastly.ErrMissingListID
	}

	path, err := common.BuildPath(i.Scope, "lists", *i.ListID)
	if err != nil {
		return nil, fmt.Errorf("failed to build API path: %w", err)
	}

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
