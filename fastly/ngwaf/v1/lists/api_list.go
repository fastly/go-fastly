package lists

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v11/fastly"
	"github.com/fastly/go-fastly/v11/fastly/ngwaf/v1/scope"
)

// ListInput specifies the information needed for the List() function to perform
// the operation.
type ListInput struct {
	// Scope defines where the list is located, including its type (e.g.,
	// "workspace" or "account") and the specific IDs it applies to (required).
	Scope *scope.Scope
}

// ListLists retrieves a list of lists for the given workspace.
func ListLists(ctx context.Context, c *fastly.Client, i *ListInput) (*Lists, error) {
	if i.Scope == nil {
		return nil, fastly.ErrMissingScope
	}

	path, err := scope.BuildPath(i.Scope, "lists", "")
	if err != nil {
		return nil, fmt.Errorf("failed to build API path: %w", err)
	}

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var lists *Lists
	if err := json.NewDecoder(resp.Body).Decode(&lists); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return lists, nil
}
