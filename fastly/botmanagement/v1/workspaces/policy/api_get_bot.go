package policy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// GetBotInput specifies the information needed to get a bot.
type GetBotInput struct {
	// BotID is the ID of the bot (required).
	BotID *string `json:"-"`
	// CategoryID is the ID of the bot category (required).
	CategoryID *string `json:"-"`
	// WorkspaceID is the ID of the bot management workspace (required).
	WorkspaceID *string `json:"-"`
}

// GetBot retrieves a specified bot.
func GetBot(ctx context.Context, c *fastly.Client, i *GetBotInput) (*Bot, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.CategoryID == nil {
		return nil, fastly.ErrMissingCategoryID
	}
	if i.BotID == nil {
		return nil, fastly.ErrMissingBotID
	}

	path := fastly.ToSafeURL("bot-management", "v1", "workspaces", *i.WorkspaceID, "policy", "categories", *i.CategoryID, "bots", *i.BotID)

	resp, err := c.GetJSON(ctx, path, fastly.CreateRequestOptions())
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
