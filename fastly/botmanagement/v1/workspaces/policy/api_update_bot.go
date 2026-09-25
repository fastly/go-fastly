package policy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// UpdateBotInput specifies the information needed to update a bot.
type UpdateBotInput struct {
	// Action is the action applied to requests from the bot. Must be one of
	// `allow`, `block`, `challenge`, or `inherit`, which applies the action of
	// the category the bot belongs to (required).
	Action *Action `json:"action"`
	// BotID is the ID of the bot (required).
	BotID *string `json:"-"`
	// CategoryID is the ID of the bot category (required).
	CategoryID *string `json:"-"`
	// WorkspaceID is the ID of the bot management workspace (required).
	WorkspaceID *string `json:"-"`
}

// UpdateBot updates a specified bot.
func UpdateBot(ctx context.Context, c *fastly.Client, i *UpdateBotInput) (*Bot, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.CategoryID == nil {
		return nil, fastly.ErrMissingCategoryID
	}
	if i.BotID == nil {
		return nil, fastly.ErrMissingBotID
	}
	if i.Action == nil {
		return nil, fastly.ErrMissingAction
	}

	path := fastly.ToSafeURL("bot-management", "v1", "workspaces", *i.WorkspaceID, "policy", "categories", *i.CategoryID, "bots", *i.BotID)

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bot *Bot
	if err := json.NewDecoder(resp.Body).Decode(&bot); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return bot, nil
}
