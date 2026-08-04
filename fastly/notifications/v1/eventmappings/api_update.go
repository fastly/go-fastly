package eventmappings

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// UpdateInput specifies the information needed for the Update() function to
// perform the operation. Update replaces the entire event mapping, so all
// fields must be provided; omitted fields are not preserved from the
// previous version.
type UpdateInput struct {
	// MappingID is the ID of the event mapping (required).
	MappingID *string `json:"-"`
	// Name is a descriptive name for the mapping (required).
	Name *string `json:"name"`
	// Description is an optional description of the mapping.
	Description *string `json:"description,omitempty"`
	// ScopeType is the category of Fastly resource the mapping applies to:
	// ScopeTypeAccount, ScopeTypeVCL, ScopeTypeWasm, or ScopeTypeNGWAF
	// (required).
	ScopeType *string `json:"scope_type"`
	// ScopeIDs is the specific service or workspace IDs to scope the mapping
	// to. Must be empty or omitted when ScopeType is ScopeTypeAccount.
	ScopeIDs []string `json:"scope_ids,omitempty"`
	// EventTypes is the audit event types that trigger a notification
	// (required).
	EventTypes []string `json:"event_types"`
	// IntegrationIDs is the IDs of the integrations that should receive
	// notifications (required, must be non-empty).
	IntegrationIDs []string `json:"integration_ids"`
}

// Update replaces an existing event mapping. The mapping is always set to
// MappingStatusActive on update.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*EventMapping, error) {
	if i.MappingID == nil {
		return nil, fastly.ErrMissingMappingID
	}
	if i.Name == nil {
		return nil, fastly.ErrMissingName
	}
	if i.ScopeType == nil {
		return nil, fastly.ErrMissingScopeType
	}
	if len(i.EventTypes) == 0 {
		return nil, fastly.ErrMissingEventTypes
	}
	if len(i.IntegrationIDs) == 0 {
		return nil, fastly.ErrMissingIntegrationIDs
	}

	path := fastly.ToSafeURL("notifications", "v1", "event-mappings", *i.MappingID)

	resp, err := c.PutJSON(ctx, path, i, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var em *EventMapping
	if err := json.NewDecoder(resp.Body).Decode(&em); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return em, nil
}
