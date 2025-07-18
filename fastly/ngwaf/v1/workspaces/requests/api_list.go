package requests

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v11/fastly"
)

// ListInput specifies the information needed for the List() function
// to perform the operation.
type ListInput struct {
	// Limit is the limit on how many results are
	// returned (required).
	Limit *int
	// Page is the page number of the collection to request
	// (Default 100).
	Page *int
	// Query is a search query string. Please read the Search
	// Syntax for help (required).
	// https://www.fastly.com/documentation/guides/next-gen-waf/reference/searching-for-requests/
	Query *string
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// List retrieves a list of requests in the specified workspace, with
// optional filtering and pagination.
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*Requests, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	requestOptions := fastly.CreateRequestOptions()
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.Page != nil {
		requestOptions.Params["page"] = strconv.Itoa(*i.Page)
	}
	if i.Query != nil {
		requestOptions.Params["q"] = *i.Query
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "requests")

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var requests *Requests
	if err := json.NewDecoder(resp.Body).Decode(&requests); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return requests, nil
}
