package rules

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// CreateInput specifies the information needed for the Create() function to
// perform the operation.
type CreateInput struct {
	// Action determines where matching requests are routed (required).
	Action *Action `json:"action"`
	// Conditions are the matching criteria evaluated against the incoming
	// request. Omit or leave empty to create the default (catch-all) rule for
	// the path.
	Conditions []Condition `json:"conditions,omitempty"`
	// PathID is the path identifier (required).
	PathID *string `json:"-"`
	// RoutingConfigID is the routing config identifier (required).
	RoutingConfigID *string `json:"-"`
}

// Create creates a new rule within the specified path.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*Data, error) {
	if i.RoutingConfigID == nil || *i.RoutingConfigID == "" {
		return nil, fastly.ErrMissingRoutingConfigID
	}
	if i.PathID == nil || *i.PathID == "" {
		return nil, fastly.ErrMissingPathID
	}
	if i.Action == nil {
		return nil, fastly.ErrMissingAction
	}

	path := fastly.ToSafeURL("domain-management", "v1", "routing-configs", *i.RoutingConfigID, "paths", *i.PathID, "rules")

	resp, err := c.PostJSON(ctx, path, i, fastly.CreateRequestOptions())
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
