package signals

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v18/fastly"
	ngwafv1 "github.com/fastly/go-fastly/v18/fastly/ngwaf/v1"
	"github.com/fastly/go-fastly/v18/fastly/ngwaf/v1/scope"
)

// ListInput specifies the information needed for the List() function
// to perform the operation.
type ListInput struct {
	// Limit how many results are returned. Defaults to the API maximum
	// of 1000.
	Limit *int
	// Scope defines where the signal is located, including its type (e.g.,
	// "workspace" or "account") and the specific IDs it applies to (required).
	Scope *scope.Scope
}

// List retrieves a list of signals, up to 1000 by default (the API's own
// default is 100).
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*Signals, error) {
	requestOptions := fastly.CreateRequestOptions()
	limit := ngwafv1.MaxLimit
	if i.Limit != nil {
		limit = *i.Limit
	}
	requestOptions.Params["limit"] = strconv.Itoa(limit)

	path, err := scope.BuildPath(i.Scope, "signals", "")
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
