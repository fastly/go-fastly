package workspaces

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v17/fastly"
)

// ListInput specifies the information needed for the List()
// function to perform the operation.
type ListInput struct {
	// Limit is how many results are returned per page.
	Limit *int
	// ServiceID is an alphanumeric string identifying the service.
	ServiceID *string
}

// List retrieves all workspaces, automatically paginating through all pages.
func List(ctx context.Context, c *fastly.Client, i *ListInput) ([]Workspace, error) {
	var (
		out    []Workspace
		cursor *string
	)
	for {
		page, err := listPage(ctx, c, i, cursor)
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

// listPage retrieves a single page of workspaces.
func listPage(ctx context.Context, c *fastly.Client, i *ListInput, cursor *string) (*Workspaces, error) {
	path := fastly.ToSafeURL("bot-management", "v1", "workspaces")

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
	if i.ServiceID != nil {
		requestOptions.Params["service_id"] = *i.ServiceID
	}

	resp, err := c.GetJSON(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var workspaces *Workspaces
	if err := json.NewDecoder(resp.Body).Decode(&workspaces); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return workspaces, nil
}
