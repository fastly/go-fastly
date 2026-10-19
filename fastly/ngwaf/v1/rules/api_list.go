package rules

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v17/fastly"
	"github.com/fastly/go-fastly/v17/fastly/ngwaf/v1/scope"
)

// ListInput specifies the information needed for the List() function
// to perform the operation.
type ListInput struct {
	// Action filter results based on action.
	Action *string
	// Enabled filter results based on enabled.
	Enabled *bool
	// Limit how many results are returned per page.
	Limit *int
	// Scope defines where the rule is applied, including its type
	// (e.g., "workspace" or "account") and the specific IDs it
	// applies to (required).
	Scope *scope.Scope
	// Types filter results based on types (accepts more than one
	// value and performs a union across rules of given types).
	Types *string
}

// List retrieves all rules, automatically paginating through every page.
func List(ctx context.Context, c *fastly.Client, i *ListInput) ([]Rule, error) {
	if i.Scope == nil {
		return nil, fastly.ErrMissingScope
	}

	out := []Rule{}
	for page := 1; ; page++ {
		r, err := listPage(ctx, c, i, page)
		if err != nil {
			return nil, err
		}
		out = append(out, r.Data...)

		// Stop on an empty page or once meta.total is reached.
		if len(r.Data) == 0 || len(out) >= r.Meta.Total {
			break
		}
	}

	return out, nil
}

// listPage retrieves a single page of rules.
func listPage(ctx context.Context, c *fastly.Client, i *ListInput, page int) (*Rules, error) {
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
	if page > 1 {
		requestOptions.Params["page"] = strconv.Itoa(page)
	}
	if i.Types != nil {
		requestOptions.Params["types"] = *i.Types
	}

	path, err := scope.BuildPath(i.Scope, "rules", "")
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
