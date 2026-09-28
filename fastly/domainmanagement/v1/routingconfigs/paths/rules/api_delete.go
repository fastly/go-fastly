package rules

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v17/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// PathID is the path identifier (required).
	PathID *string
	// RoutingConfigID is the routing config identifier (required).
	RoutingConfigID *string
	// RuleID of the rule to delete (required).
	RuleID *string
}

// Delete deletes the specified rule.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.RoutingConfigID == nil || *i.RoutingConfigID == "" {
		return fastly.ErrMissingRoutingConfigID
	}
	if i.PathID == nil || *i.PathID == "" {
		return fastly.ErrMissingPathID
	}
	if i.RuleID == nil || *i.RuleID == "" {
		return fastly.ErrMissingRuleID
	}

	path := fastly.ToSafeURL("domain-management", "v1", "routing-configs", *i.RoutingConfigID, "paths", *i.PathID, "rules", *i.RuleID)

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
