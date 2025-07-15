package rules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fastly/go-fastly/v10/fastly"
	"github.com/fastly/go-fastly/v10/fastly/ngwaf/v1/common"
)

// DeleteInput specifies the information needed for the Delete()
// function to perform the operation.
type DeleteInput struct {
	// RuleID is the rule identifier (required).
	RuleID *string
	// Scope defines where the rule is applied, including its type (e.g., "workspace" or "account") and the specific IDs it applies to (required).
	Scope *common.Scope
}

// Delete deletes the specified rule.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.RuleID == nil {
		return fastly.ErrMissingRuleID
	}
	if i.Scope == nil {
		return fastly.ErrMissingScope
	}

	path, err := common.BuildPath(i.Scope, "rules", *i.RuleID)
	if err != nil {
		return fmt.Errorf("failed to build API path: %w", err)
	}

	resp, err := c.Delete(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fastly.NewHTTPError(resp)
	}

	return nil
}
