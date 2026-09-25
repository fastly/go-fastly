package policy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v17/fastly"
	"github.com/fastly/go-fastly/v17/fastly/botmanagement/v1/workspaces"
)

func TestClient_Policy(t *testing.T) {
	ctx := context.TODO()

	var err error

	// Find a workspace to test against. Workspaces cannot be created through
	// the API, so the test account must already have at least one.
	var ws []workspaces.Workspace
	fastly.Record(t, "list_workspaces", func(c *fastly.Client) {
		ws, err = workspaces.List(ctx, c, &workspaces.ListInput{})
	})
	require.NoError(t, err)
	require.NotEmpty(t, ws, "expected at least one workspace")
	workspaceID := ws[0].ID

	// List categories.
	var categories []Category
	fastly.Record(t, "list_categories", func(c *fastly.Client) {
		categories, err = ListCategories(ctx, c, &ListCategoriesInput{
			WorkspaceID: new(workspaceID),
		})
	})
	require.NoError(t, err)
	require.NotEmpty(t, categories, "expected at least one category")

	// List categories 10 per page to exercise auto-pagination. The result
	// should match the single-page listing.
	var paginatedCategories []Category
	fastly.Record(t, "list_categories_paginated", func(c *fastly.Client) {
		paginatedCategories, err = ListCategories(ctx, c, &ListCategoriesInput{
			Limit:       new(10),
			WorkspaceID: new(workspaceID),
		})
	})
	require.NoError(t, err)
	require.Len(t, paginatedCategories, len(categories))

	// Use a category with enough bots to span multiple pages when listing
	// its bots 10 per page. Categories are defined by Fastly, so the ID is
	// stable across accounts.
	const testCategoryID = "cat_monitoring_and_site_tools"

	var category Category
	for _, cat := range categories {
		if cat.CategoryID == testCategoryID {
			category = cat
			break
		}
	}
	require.Equal(t, testCategoryID, category.CategoryID, "expected category %q to exist", testCategoryID)

	// Get the category.
	var gotCategory *Category
	fastly.Record(t, "get_category", func(c *fastly.Client) {
		gotCategory, err = GetCategory(ctx, c, &GetCategoryInput{
			CategoryID:  new(category.CategoryID),
			WorkspaceID: new(workspaceID),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, gotCategory)
	require.Equal(t, category.CategoryID, gotCategory.CategoryID)
	require.Equal(t, category.Name, gotCategory.Name)
	require.Equal(t, category.Action.Type, gotCategory.Action.Type)

	// Restore the original category action at the end.
	defer func() {
		fastly.Record(t, "restore_category", func(c *fastly.Client) {
			_, err = UpdateCategory(ctx, c, &UpdateCategoryInput{
				Action:      &Action{Type: category.Action.Type},
				CategoryID:  new(category.CategoryID),
				WorkspaceID: new(workspaceID),
			})
		})
		if err != nil {
			t.Errorf("error restoring category: %v", err)
		}
	}()

	// Update the category action to a different value.
	updatedCategoryAction := "block"
	if category.Action.Type == updatedCategoryAction {
		updatedCategoryAction = "challenge"
	}

	var updatedCategory *Category
	fastly.Record(t, "update_category", func(c *fastly.Client) {
		updatedCategory, err = UpdateCategory(ctx, c, &UpdateCategoryInput{
			Action:      &Action{Type: updatedCategoryAction},
			CategoryID:  new(category.CategoryID),
			WorkspaceID: new(workspaceID),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, updatedCategory)
	require.Equal(t, category.CategoryID, updatedCategory.CategoryID)
	require.Equal(t, updatedCategoryAction, updatedCategory.Action.Type)

	// List all bots in the workspace. Each bot should report its category.
	var workspaceBots []Bot
	fastly.Record(t, "list_bots", func(c *fastly.Client) {
		workspaceBots, err = ListBots(ctx, c, &ListBotsInput{
			WorkspaceID: new(workspaceID),
		})
	})
	require.NoError(t, err)
	require.NotEmpty(t, workspaceBots, "expected at least one bot")
	for _, b := range workspaceBots {
		require.NotEmpty(t, b.CategoryID, "expected bot %q to have a category", b.BotID)
	}

	// List the bots in the category.
	var categoryBots []Bot
	fastly.Record(t, "list_category_bots", func(c *fastly.Client) {
		categoryBots, err = ListBots(ctx, c, &ListBotsInput{
			CategoryID:  new(category.CategoryID),
			WorkspaceID: new(workspaceID),
		})
	})
	require.NoError(t, err)
	require.NotEmpty(t, categoryBots, "expected at least one bot in category %q", category.CategoryID)

	// List the bots in the category 10 per page to exercise
	// auto-pagination. The result should match the single-page listing.
	var paginatedCategoryBots []Bot
	fastly.Record(t, "list_category_bots_paginated", func(c *fastly.Client) {
		paginatedCategoryBots, err = ListBots(ctx, c, &ListBotsInput{
			CategoryID:  new(category.CategoryID),
			Limit:       new(10),
			WorkspaceID: new(workspaceID),
		})
	})
	require.NoError(t, err)
	require.Len(t, paginatedCategoryBots, len(categoryBots))

	bot := categoryBots[0]

	// Get the bot.
	var gotBot *Bot
	fastly.Record(t, "get_bot", func(c *fastly.Client) {
		gotBot, err = GetBot(ctx, c, &GetBotInput{
			BotID:       new(bot.BotID),
			CategoryID:  new(category.CategoryID),
			WorkspaceID: new(workspaceID),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, gotBot)
	require.Equal(t, bot.BotID, gotBot.BotID)
	require.Equal(t, bot.Name, gotBot.Name)
	require.Equal(t, bot.Action.Type, gotBot.Action.Type)

	// Restore the original bot action at the end.
	defer func() {
		fastly.Record(t, "restore_bot", func(c *fastly.Client) {
			_, err = UpdateBot(ctx, c, &UpdateBotInput{
				Action:      &Action{Type: bot.Action.Type},
				BotID:       new(bot.BotID),
				CategoryID:  new(category.CategoryID),
				WorkspaceID: new(workspaceID),
			})
		})
		if err != nil {
			t.Errorf("error restoring bot: %v", err)
		}
	}()

	// Update the bot action to a different value.
	updatedBotAction := "inherit"
	if bot.Action.Type == updatedBotAction {
		updatedBotAction = "allow"
	}

	var updatedBot *Bot
	fastly.Record(t, "update_bot", func(c *fastly.Client) {
		updatedBot, err = UpdateBot(ctx, c, &UpdateBotInput{
			Action:      &Action{Type: updatedBotAction},
			BotID:       new(bot.BotID),
			CategoryID:  new(category.CategoryID),
			WorkspaceID: new(workspaceID),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, updatedBot)
	require.Equal(t, bot.BotID, updatedBot.BotID)
	require.Equal(t, updatedBotAction, updatedBot.Action.Type)
}

func TestClient_Policy_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := ListCategories(ctx, fastly.TestClient, &ListCategoriesInput{
		WorkspaceID: nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingWorkspaceID)

	_, err = GetCategory(ctx, fastly.TestClient, &GetCategoryInput{
		WorkspaceID: nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingWorkspaceID)

	_, err = GetCategory(ctx, fastly.TestClient, &GetCategoryInput{
		CategoryID:  nil,
		WorkspaceID: new("workspace"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingCategoryID)

	_, err = UpdateCategory(ctx, fastly.TestClient, &UpdateCategoryInput{
		WorkspaceID: nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingWorkspaceID)

	_, err = UpdateCategory(ctx, fastly.TestClient, &UpdateCategoryInput{
		CategoryID:  nil,
		WorkspaceID: new("workspace"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingCategoryID)

	_, err = UpdateCategory(ctx, fastly.TestClient, &UpdateCategoryInput{
		Action:      nil,
		CategoryID:  new("category"),
		WorkspaceID: new("workspace"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingAction)

	_, err = ListBots(ctx, fastly.TestClient, &ListBotsInput{
		WorkspaceID: nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingWorkspaceID)

	_, err = GetBot(ctx, fastly.TestClient, &GetBotInput{
		WorkspaceID: nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingWorkspaceID)

	_, err = GetBot(ctx, fastly.TestClient, &GetBotInput{
		CategoryID:  nil,
		WorkspaceID: new("workspace"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingCategoryID)

	_, err = GetBot(ctx, fastly.TestClient, &GetBotInput{
		BotID:       nil,
		CategoryID:  new("category"),
		WorkspaceID: new("workspace"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingBotID)

	_, err = UpdateBot(ctx, fastly.TestClient, &UpdateBotInput{
		WorkspaceID: nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingWorkspaceID)

	_, err = UpdateBot(ctx, fastly.TestClient, &UpdateBotInput{
		CategoryID:  nil,
		WorkspaceID: new("workspace"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingCategoryID)

	_, err = UpdateBot(ctx, fastly.TestClient, &UpdateBotInput{
		BotID:       nil,
		CategoryID:  new("category"),
		WorkspaceID: new("workspace"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingBotID)

	_, err = UpdateBot(ctx, fastly.TestClient, &UpdateBotInput{
		Action:      nil,
		BotID:       new("bot"),
		CategoryID:  new("category"),
		WorkspaceID: new("workspace"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingAction)
}
