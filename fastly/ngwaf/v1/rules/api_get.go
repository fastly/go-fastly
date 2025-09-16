package rules

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v11/fastly"
	"github.com/fastly/go-fastly/v11/fastly/ngwaf/v1/scope"
)

// GetInput specifies the information needed for the Get() function
// to perform the operation.
type GetInput struct {
	// RuleID is the rule identifier (required).
	RuleID *string
	// Scope defines where the rule is applied, including its type
	// (e.g., "workspace" or "account") and the specific IDs it
	// applies to (required).
	Scope *scope.Scope
}

// Get retrieves the specified rule.
func Get(ctx context.Context, c *fastly.Client, i *GetInput) (*Rule, error) {
	if i.RuleID == nil {
		return nil, fastly.ErrMissingRuleID
	}
	if i.Scope == nil {
		return nil, fastly.ErrMissingScope
	}

	path, err := scope.BuildPath(i.Scope, "rules", *i.RuleID)
	if err != nil {
		return nil, fmt.Errorf("failed to build API path: %w", err)
	}

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r *Rule
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return r, nil
}
