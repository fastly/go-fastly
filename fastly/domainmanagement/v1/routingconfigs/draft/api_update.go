package draft

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// UpdateInput specifies the information needed for the Update() function to
// perform the operation.
type UpdateInput struct {
	// Comment is a descriptive note about the changes made in the draft
	// (required).
	Comment *string `json:"comment"`
	// RoutingConfigID is the routing config identifier (required).
	RoutingConfigID *string `json:"-"`
}

// Update sets the comment on the specified routing config's draft version.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*Data, error) {
	if i.RoutingConfigID == nil || *i.RoutingConfigID == "" {
		return nil, fastly.ErrMissingRoutingConfigID
	}
	if i.Comment == nil {
		return nil, fastly.ErrMissingComment
	}

	path := fastly.ToSafeURL("domain-management", "v1", "routing-configs", *i.RoutingConfigID, "draft")

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
