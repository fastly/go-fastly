package policy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v19/fastly"
)

// UpdateCategoryBotsInput specifies the information needed to update every
// bot in a category.
type UpdateCategoryBotsInput struct {
	// Action is the action to apply to every bot in the category (required).
	//
	// The only supported value is `inherit`, which removes each bot's
	// own action, so it uses the action of the category.
	Action *Action `json:"action"`
	// CategoryID is the ID of the bot category (required).
	CategoryID *string `json:"-"`
	// WorkspaceID is the ID of the bot management workspace (required).
	WorkspaceID *string `json:"-"`
}

// UpdateCategoryBots sets the action for every bot in a specified category.
// The category's own action is unchanged. The update applies to every bot
// in the category or to none, and repeating it gives the same result.
//
// It returns every bot in the category after the update, including bots
// that already inherited the category's action.
func UpdateCategoryBots(ctx context.Context, c *fastly.Client, i *UpdateCategoryBotsInput) ([]Bot, error) {
	if fastly.IsNilOrEmpty(i.WorkspaceID) {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if fastly.IsNilOrEmpty(i.CategoryID) {
		return nil, fastly.ErrMissingCategoryID
	}
	if i.Action == nil {
		return nil, fastly.ErrMissingAction
	}

	path := fastly.ToSafeURL("bot-management", "v1", "workspaces", *i.WorkspaceID, "policy", "categories", *i.CategoryID, "bots")

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bots *Bots
	if err := json.NewDecoder(resp.Body).Decode(&bots); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}
	if bots == nil {
		return nil, nil
	}

	return bots.Data, nil
}
