package thresholds

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v18/fastly"
)

// ListInput specifies the information needed for the List() function
// to perform the operation.
type ListInput struct {
	// Limit how many results are returned. Defaults to the API maximum
	// of 1000.
	Limit *int
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// maxLimit is the largest page size the API accepts.
const maxLimit = 1000

// List retrieves a list of thresholds, up to 1000 by default (the API's
// own default is 100).
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*Thresholds, error) {
	requestOptions := fastly.CreateRequestOptions()
	limit := maxLimit
	if i.Limit != nil {
		limit = *i.Limit
	}
	requestOptions.Params["limit"] = strconv.Itoa(limit)

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "thresholds")

	resp, err := c.Get(ctx, path, requestOptions)
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
