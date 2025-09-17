package lists

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v11/fastly"
	"github.com/fastly/go-fastly/v11/fastly/ngwaf/v1/scope"
)

// GetInput specifies the information needed for the Get() function to
// perform the operation.
type GetInput struct {
	// ListID is the workspace identifier (required).
	ListID *string
	// Scope defines where the list is applied, including its type (e.g.,
	// "workspace" or "account") and the specific IDs it applies to (required).
	Scope *scope.Scope
}

// Get retrieves the specified list.
func Get(ctx context.Context, c *fastly.Client, i *GetInput) (*List, error) {
	if i.ListID == nil {
		return nil, fastly.ErrMissingListID
	}
	if i.Scope == nil {
		return nil, fastly.ErrMissingScope
	}

	path, err := scope.BuildPath(i.Scope, "lists", *i.ListID)
	if err != nil {
		return nil, fmt.Errorf("failed to build API path: %w", err)
	}

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
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
