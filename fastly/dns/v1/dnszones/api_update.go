package dnszones

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v14/fastly"
)

// UpdateInput specifies the information needed to update a zone.
type UpdateInput struct {
	// ZoneID is the Zone Identifier (UUID) (required).
	ZoneID *string `json:"-"`
	// Description is a freeform descriptive note.
	Description *fastly.Nullable[string] `json:"description,omitempty"`
	// XfrConfigInbound contains attributes associated with inbound zone transfers.
	XfrConfigInbound *XfrConfigInboundInput `json:"xfr_config_inbound,omitempty"`
}

// Update updates an existing DNS Zone.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*Zone, error) {
	if i.ZoneID == nil {
		return nil, fastly.ErrMissingID
	}

	path := fastly.ToSafeURL("dns", "v1", "zones", *i.ZoneID)

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var zone *Zone
	if err := json.NewDecoder(resp.Body).Decode(&zone); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return zone, nil
}
