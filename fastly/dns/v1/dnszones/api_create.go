package dnszones

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v14/fastly"
)

// CreateInput specifies the information needed to create a zone.
type CreateInput struct {
	// Name is the domain name for the DNS zone (required).
	Name *string `json:"name"`
	// Type is the type of the DNS zone (required).
	Type *string `json:"type"`
	// Description is a freeform descriptive note.
	Description *string `json:"description,omitempty"`
	// XfrConfigInbound contains attributes associated with inbound zone transfers.
	XfrConfigInbound *XfrConfigInboundInput `json:"xfr_config_inbound,omitempty"`
}

// Create creates a new DNS Zone.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*Zone, error) {
	if i.Name == nil {
		return nil, fastly.ErrMissingName
	}
	if i.Type == nil {
		return nil, fastly.ErrMissingType
	}

	resp, err := c.PostJSON(ctx, "/dns/v1/zones", i, fastly.CreateRequestOptions())
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
