package rules

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v18/fastly"
	"github.com/fastly/go-fastly/v18/fastly/ngwaf/v1/scope"
)

// ListInput specifies the information needed for the List() function
// to perform the operation.
type ListInput struct {
	// Action filter results based on action.
	Action *string
	// Enabled filter results based on enabled.
	Enabled *bool
	// Limit is the maximum number of results to return.
	Limit *int
	// Scope defines where the rule is applied, including its type
	// (e.g., "workspace" or "account") and the specific IDs it
	// applies to (required).
	Scope *scope.Scope
	// Types filter results based on types (accepts more than one
	// value and performs a union across rules of given types).
	Types *string
}

// maxPageSize is the largest page size the API accepts.
const maxPageSize = 1000

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

		// Stop on an empty page, once Limit results have been
		// collected, or once meta.total is reached.
		if i.Limit != nil && *i.Limit > 0 && len(out) >= *i.Limit {
			out = out[:*i.Limit]
			break
		}
		if len(r.Data) == 0 || len(out) >= r.Meta.Total {
			break
		}
	}

	return out, nil
}

// listPage retrieves a single page of rules.
func listPage(ctx context.Context, c *fastly.Client, i *ListInput, page int) (*Rules, error) {
	// Request the largest page the API allows, unless Limit asks for fewer.
	perPage := maxPageSize
	if i.Limit != nil && *i.Limit > 0 && *i.Limit < perPage {
		perPage = *i.Limit
	}

	requestOptions := fastly.CreateRequestOptions()
	requestOptions.Params["limit"] = strconv.Itoa(perPage)
	if i.Action != nil {
		requestOptions.Params["action"] = *i.Action
	}
	if i.Enabled != nil {
		requestOptions.Params["enabled"] = strconv.FormatBool(*i.Enabled)
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
