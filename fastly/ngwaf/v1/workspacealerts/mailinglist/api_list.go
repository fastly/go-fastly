package mailinglist

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
	// Limit how many results are returned (optional).
	Limit *int
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// List retrieves a list of workspaces alerts.
func List(c *fastly.Client, i *ListInput) (*WorkspaceAlerts, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	requestOptions := fastly.CreateRequestOptions(i.Context)
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "alerts")

	resp, err := c.Get(path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var was *WorkspaceAlerts
	if err := json.NewDecoder(resp.Body).Decode(&was); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return was, nil
}