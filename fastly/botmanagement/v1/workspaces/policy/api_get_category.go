package policy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// GetCategoryInput specifies the information needed to get a category.
type GetCategoryInput struct {
	// CategoryID is the ID of the bot category (required).
	CategoryID *string `json:"-"`
	// WorkspaceID is the ID of the bot management workspace (required).
	WorkspaceID *string `json:"-"`
}

// GetCategory retrieves a specified category.
func GetCategory(ctx context.Context, c *fastly.Client, i *GetCategoryInput) (*Category, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.CategoryID == nil {
		return nil, fastly.ErrMissingCategoryID
	}

	path := fastly.ToSafeURL("bot-management", "v1", "workspaces", *i.WorkspaceID, "policy", "categories", *i.CategoryID)

	resp, err := c.GetJSON(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var category *Category
	if err := json.NewDecoder(resp.Body).Decode(&category); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return category, nil
}
