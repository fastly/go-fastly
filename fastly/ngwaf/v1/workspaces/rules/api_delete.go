package rules

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v10/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// RuleID is the rule identifier (required).
	RuleID *string
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Delete deletes the specified rule.
func Delete(c *fastly.Client, i *DeleteInput) error {
	if i.WorkspaceID == nil {
		return fastly.ErrMissingWorkspaceID
	}
	if i.RuleID == nil {
		return fastly.ErrMissingRuleID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "rules", *i.RuleID)

	resp, err := c.Delete(path, fastly.CreateRequestOptions(i.Context))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fastly.NewHTTPError(resp)
	}

	return nil
}
