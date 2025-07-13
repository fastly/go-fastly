package rules

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v10/fastly"
)

// DeleteInput specifies the information needed for the Delete()
// function to perform the operation.
type DeleteInput struct {
	// RuleID is the rule identifier (required).
	RuleID *string
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Delete deletes the specified rule.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.WorkspaceID == nil {
		return fastly.ErrMissingWorkspaceID
	}
	if i.RuleID == nil {
		return fastly.ErrMissingRuleID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "rules", *i.RuleID)

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
