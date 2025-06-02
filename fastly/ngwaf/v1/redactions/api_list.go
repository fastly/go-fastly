package redactions

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
	// Limit how many results are returned.
	Limit *int
	// Mode filter results based on mode.
	Mode *string
	// Page number of the collection to request.
	Page *int
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// List retrieves a list of workspaces, with optional filtering and pagination.
func List(c *fastly.Client, i *ListInput) (*Redactions, error) {
	requestOptions := fastly.CreateRequestOptions(i.Context)
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.Mode != nil {
		requestOptions.Params["mode"] = *i.Mode
	}
	if i.Page != nil {
		requestOptions.Params["page"] = strconv.Itoa(*i.Page)
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "redactions")

	resp, err := c.Get(path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var redactions *Redactions
	if err := json.NewDecoder(resp.Body).Decode(&redactions); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return redactions, nil
}
