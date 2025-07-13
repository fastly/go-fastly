package thresholds

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v11/fastly"
)

// ListInput specifies the information needed for the List() function to perform
// the operation.
type ListInput struct {
	// Limit how many results are returned.
	Limit *int
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// List retrieves a list of thresholds, with optional limiting.
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*Thresholds, error) {
	requestOptions := fastly.CreateRequestOptions()
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}

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
