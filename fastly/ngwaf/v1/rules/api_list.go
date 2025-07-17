package rules

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
	// Action filter results based on action.
	Action *string
	// Enabled filter results based on enabled.
	Enabled *bool
	// Limit how many results are returned.
	Limit *int
	// Page number of the collection to request.
	Page *int
	// Scope defines where the rule is applied, including its type
	// (e.g., "workspace" or "account") and the specific IDs it
	// applies to (required).
	Scope *common.Scope
	// Types filter results based on types (accepts more than one
	// value and performs a union across rules of given types).
	Types *string
}

// List retrieves a list of rules, with optional filtering and
// pagination.
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*Rules, error) {
	if i.Scope == nil {
		return nil, fastly.ErrMissingScope
	}

	requestOptions := fastly.CreateRequestOptions()
	if i.Action != nil {
		requestOptions.Params["action"] = *i.Action
	}
	if i.Enabled != nil {
		requestOptions.Params["enabled"] = strconv.FormatBool(*i.Enabled)
	}
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.Page != nil {
		requestOptions.Params["page"] = strconv.Itoa(*i.Page)
	}
	if i.Types != nil {
		requestOptions.Params["types"] = *i.Types
	}

	path, err := common.BuildPath(i.Scope, "rules", "")
	if err != nil {
		return nil, fmt.Errorf("failed to build API path: %w", err)
	}

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r *Rules
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return r, nil
}
