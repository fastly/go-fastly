package rules

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// UpdateInput specifies the information needed for the Update() function to
// perform the operation.
type UpdateInput struct {
	// Action determines where matching requests are routed.
	Action *Action `json:"action,omitempty"`
	// Conditions are the matching criteria evaluated against the incoming
	// request. Pass a pointer to an empty slice to turn the rule into the
	// default (catch-all) rule for its path.
	Conditions *[]Condition `json:"conditions,omitempty"`
	// PathID is the path identifier (required).
	PathID *string `json:"-"`
	// RoutingConfigID is the routing config identifier (required).
	RoutingConfigID *string `json:"-"`
	// RuleID is the rule identifier (required).
	RuleID *string `json:"-"`
}

// Update updates the specified rule.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*Data, error) {
	if i.RoutingConfigID == nil || *i.RoutingConfigID == "" {
		return nil, fastly.ErrMissingRoutingConfigID
	}
	if i.PathID == nil || *i.PathID == "" {
		return nil, fastly.ErrMissingPathID
	}
	if i.RuleID == nil || *i.RuleID == "" {
		return nil, fastly.ErrMissingRuleID
	}

	path := fastly.ToSafeURL("domain-management", "v1", "routing-configs", *i.RoutingConfigID, "paths", *i.PathID, "rules", *i.RuleID)

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *Data
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}
	return d, nil
}
