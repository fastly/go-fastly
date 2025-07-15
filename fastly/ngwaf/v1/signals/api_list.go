package signals

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v10/fastly"
	"github.com/fastly/go-fastly/v10/fastly/ngwaf/v1/common"
)

// ListInput specifies the information needed for the List() function
// to perform the operation.
type ListInput struct {
	// Limit how many results are returned.
	Limit *int
	// Scope defines where the signal is located, including its type (e.g.,
	// "workspace" or "account") and the specific IDs it applies to (required).
	Scope *common.Scope
}

// List retrieves a list of signals for the given workspace, with
// optional pagination.
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*Signals, error) {
	requestOptions := fastly.CreateRequestOptions()
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}

	path, err := common.BuildPath(i.Scope, "signals", "")
	if err != nil {
		return nil, fmt.Errorf("failed to build API path: %w", err)
	}

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
