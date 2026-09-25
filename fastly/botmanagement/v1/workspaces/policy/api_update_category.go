package policy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// UpdateCategoryInput specifies the information needed to update a category.
type UpdateCategoryInput struct {
	// Action is the action applied to every bot in the category that is
	// configured with the `inherit` action. Must be one of `allow`,
	// `block`, or `challenge` (required).
	Action *Action `json:"action"`
	// CategoryID is the ID of the bot category (required).
	CategoryID *string `json:"-"`
	// WorkspaceID is the ID of the bot management workspace (required).
	WorkspaceID *string `json:"-"`
}

// UpdateCategory updates a specified category.
func UpdateCategory(ctx context.Context, c *fastly.Client, i *UpdateCategoryInput) (*Category, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.CategoryID == nil {
		return nil, fastly.ErrMissingCategoryID
	}
	if i.Action == nil {
		return nil, fastly.ErrMissingAction
	}

	path := fastly.ToSafeURL("bot-management", "v1", "workspaces", *i.WorkspaceID, "policy", "categories", *i.CategoryID)

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
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
