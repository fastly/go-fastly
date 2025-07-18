package lists

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
	"github.com/fastly/go-fastly/v10/fastly/ngwaf/v1/common"
)

// CreateInput specifies the information needed for the Create()
// function to perform the operation.
type CreateInput struct {
	// Description is the description of the list.
	Description *string `json:"description,omitempty"`
	// Entries are the entries of the list (required).
	Entries *[]string `json:"entries"`
	// Name is the name of the list (required).
	Name *string `json:"name"`
	// Scope defines where the list is located, including its type (e.g.,
	// "workspace" or "account") and the specific IDs it applies to (required).
	Scope *common.Scope `json:"-"`
	// Type is the type of the list. Must be one of `string` |
	// `wildcard` | `ip` | `country` | `signal` (required).
	Type *string `json:"type"`
}

// Create creates a new list.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*List, error) {
	if i.Entries == nil {
		return nil, fastly.ErrMissingEntries
	}
	if i.Name == nil {
		return nil, fastly.ErrMissingName
	}
	if i.Type == nil {
		return nil, fastly.ErrMissingType
	}
	if i.Scope == nil {
		return nil, fastly.ErrMissingScope
	}

	path, err := common.BuildPath(i.Scope, "lists", "")
	if err != nil {
		return nil, fmt.Errorf("failed to build API path: %w", err)
	}

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
