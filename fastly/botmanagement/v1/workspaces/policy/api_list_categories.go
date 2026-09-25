package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v17/fastly"
)

// ListCategoriesInput specifies the information needed for the
// ListCategories() function to perform the operation.
type ListCategoriesInput struct {
	// Limit is how many results are returned per page.
	Limit *int
	// WorkspaceID is the ID of the bot management workspace (required).
	WorkspaceID *string
}

// ListCategories retrieves all categories in a workspace, automatically paginating through all pages.
func ListCategories(ctx context.Context, c *fastly.Client, i *ListCategoriesInput) ([]Category, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	var (
		out    []Category
		cursor *string
	)
	for {
		page, err := listCategoriesPage(ctx, c, i, cursor)
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

// listCategoriesPage retrieves a single page of categories in a workspace.
func listCategoriesPage(ctx context.Context, c *fastly.Client, i *ListCategoriesInput, cursor *string) (*Categories, error) {
	path := fastly.ToSafeURL("bot-management", "v1", "workspaces", *i.WorkspaceID, "policy", "categories")

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

	var categories *Categories
	if err := json.NewDecoder(resp.Body).Decode(&categories); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return categories, nil
}
