package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v17/fastly"
)

// ListBotsInput specifies the information needed for the ListBots()
// function to perform the operation.
type ListBotsInput struct {
	// CategoryID is the ID of the bot category. When set, only the bots in
	// that category are returned. Otherwise, all bots in the workspace are
	// returned.
	CategoryID *string
	// Limit is how many results are returned per page.
	Limit *int
	// WorkspaceID is the ID of the bot management workspace (required).
	WorkspaceID *string
}

// ListBots retrieves all bots in a workspace, or in a single category if
// CategoryID is set, automatically paginating through all pages.
func ListBots(ctx context.Context, c *fastly.Client, i *ListBotsInput) ([]Bot, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	var (
		out    []Bot
		cursor *string
	)
	for {
		page, err := listBotsPage(ctx, c, i, cursor)
		if err != nil {
			return nil, err
		}
		out = append(out, page.Data...)
		if page.Meta.NextCursor == nil || *page.Meta.NextCursor == "" {
			break
		}
		cursor = page.Meta.NextCursor
	}
	return out, nil
}

// listBotsPage retrieves a single page of bots in a workspace.
func listBotsPage(ctx context.Context, c *fastly.Client, i *ListBotsInput, cursor *string) (*Bots, error) {
	path := fastly.ToSafeURL("bot-management", "v1", "workspaces", *i.WorkspaceID, "policy", "bots")
	// Also handles listing bots in a category if CategoryID is provided.
	if i.CategoryID != nil {
		path = fastly.ToSafeURL("bot-management", "v1", "workspaces", *i.WorkspaceID, "policy", "categories", *i.CategoryID, "bots")
	}

	// limit defaults to the maximum value to optimize the number of requests made.
	limit := 1000
	if i.Limit != nil {
		limit = *i.Limit
	}

	requestOptions := fastly.CreateRequestOptions()
	if cursor != nil {
		requestOptions.Params["cursor"] = *cursor
	}
	// Always send limit so we override the default value of 100.
	requestOptions.Params["limit"] = strconv.Itoa(limit)

	resp, err := c.GetJSON(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bots *Bots
	if err := json.NewDecoder(resp.Body).Decode(&bots); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return bots, nil
}
