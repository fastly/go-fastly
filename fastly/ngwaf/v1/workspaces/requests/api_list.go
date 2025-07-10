package requests

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v10/fastly"
)

// ListInput specifies the information needed for the List() function to perform
// the operation.
type ListInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// Limit is the limit on how many results are returned. Required.
	Limit *int
	// Page is the page number of the collection to request (Default 100). Optional.
	Page *int
	// Query is a search query string. Please read the Search Syntax for help. Required.
	// https://www.fastly.com/documentation/guides/next-gen-waf/reference/searching-for-requests/
	Query *string
	// WorkspaceID is the workspace identifier (UUID).
	WorkspaceID *string
}

// List retrieves all requests in the specified workspace.
func List(c *fastly.Client, i *ListInput) (*Requests, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	requestOptions := fastly.CreateRequestOptions(i.Context)
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

	resp, err := c.Get(path, requestOptions)
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
