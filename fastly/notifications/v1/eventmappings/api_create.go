package eventmappings

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// CreateInput specifies the information needed for the Create() function to
// perform the operation.
type CreateInput struct {
	// Name is a descriptive name for the mapping (required).
	Name *string `json:"name"`
	// Description is an optional description of the mapping.
	Description *string `json:"description,omitempty"`
	// ScopeType is the category of Fastly resource the mapping applies to:
	// ScopeTypeAccount, ScopeTypeVCL, ScopeTypeWasm, or ScopeTypeNGWAF
	// (required).
	ScopeType *string `json:"scope_type"`
	// ScopeIDs is the specific service or workspace IDs to scope the mapping
	// to. Omit or provide an empty slice to apply the mapping to all
	// resources of the given scope type. Must be empty or omitted when
	// ScopeType is ScopeTypeAccount.
	ScopeIDs []string `json:"scope_ids,omitempty"`
	// EventTypes is the audit event types that trigger a notification
	// (required). Each event type must be valid for the given ScopeType.
	EventTypes []string `json:"event_types"`
	// IntegrationIDs is the IDs of the integrations that should receive
	// notifications (required). Must reference integrations belonging to
	// the account linked to the supplied token.
	IntegrationIDs []string `json:"integration_ids"`
}

// Create creates a new audit log event mapping. The mapping becomes active
// immediately.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*EventMapping, error) {
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

	path := fastly.ToSafeURL("notifications", "v1", "event-mappings")

	resp, err := c.PostJSON(ctx, path, i, fastly.CreateRequestOptions())
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
