package draft

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// GetDiffInput specifies the information needed for the GetDiff() function to
// perform the operation.
type GetDiffInput struct {
	// RoutingConfigID is the routing config identifier (required).
	RoutingConfigID *string
}

// GetDiff retrieves the differences between the active and draft versions of
// the specified routing config. The routing config must have both an active
// and a draft version.
func GetDiff(ctx context.Context, c *fastly.Client, i *GetDiffInput) (*Diff, error) {
	if i.RoutingConfigID == nil {
		return nil, fastly.ErrMissingRoutingConfigID
	}

	path := fastly.ToSafeURL("domain-management", "v1", "routing-configs", *i.RoutingConfigID, "draft", "diff")

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *Diff
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return d, nil
}
