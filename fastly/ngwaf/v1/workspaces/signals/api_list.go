package signals

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v10/fastly"
)

// ListInput specifies the information needed for the List() function
// to perform the operation.
type ListInput struct {
	// Limit how many results are returned.
	Limit *int
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// List retrieves a list of signals for the given workspace, with
// optional pagination.
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*Signals, error) {
	requestOptions := fastly.CreateRequestOptions()
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "signals")

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var signals *Signals
	if err := json.NewDecoder(resp.Body).Decode(&signals); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return signals, nil
}
