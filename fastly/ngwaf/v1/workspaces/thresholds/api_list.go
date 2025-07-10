package thresholds

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
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// List retrieves a list of thresholds, with optional limiting.
func List(c *fastly.Client, i *ListInput) (*Thresholds, error) {
	requestOptions := fastly.CreateRequestOptions(i.Context)
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "thresholds")

	resp, err := c.Get(path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var thresholds *Thresholds
	if err := json.NewDecoder(resp.Body).Decode(&thresholds); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return thresholds, nil
}
