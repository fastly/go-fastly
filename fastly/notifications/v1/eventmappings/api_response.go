package eventmappings

import "time"

const (
	// ScopeTypeAccount scopes a mapping account-wide, regardless of any
	// specific service or workspace.
	ScopeTypeAccount = "account"
	// ScopeTypeVCL scopes a mapping to specific CDN (VCL) services.
	ScopeTypeVCL = "vcl"
	// ScopeTypeWasm scopes a mapping to specific Compute (Wasm) services.
	ScopeTypeWasm = "wasm"
	// ScopeTypeNGWAF scopes a mapping to specific Next-Gen WAF workspaces.
	ScopeTypeNGWAF = "ngwaf"
)

const (
	// MappingStatusActive indicates the mapping has at least one integration
	// ID and is permitted to send notifications.
	MappingStatusActive = "active"
	// MappingStatusInactive indicates the mapping has no integration IDs and
	// will not send notifications.
	MappingStatusInactive = "inactive"
)

// EventMapping is the API response structure for the create, get, update,
// and list operations.
type EventMapping struct {
	// ID is the unique identifier for the mapping.
	ID string `json:"id"`
	// CustomerID is the ID of the customer this mapping belongs to.
	CustomerID string `json:"customer_id"`
	// Name is the descriptive name for the mapping.
	Name string `json:"name"`
	// Description is the description of the mapping.
	Description string `json:"description,omitempty"`
	// ScopeType is the category of Fastly resource the mapping applies to.
	// One of ScopeTypeAccount, ScopeTypeVCL, ScopeTypeWasm, or
	// ScopeTypeNGWAF.
	ScopeType string `json:"scope_type"`
	// ScopeIDs is the specific service or workspace IDs the mapping is
	// scoped to. Empty when the mapping applies to all resources of the
	// given scope type.
	ScopeIDs []string `json:"scope_ids"`
	// EventTypes is the audit event types that trigger a notification.
	EventTypes []string `json:"event_types"`
	// IntegrationIDs is the IDs of the integrations that receive
	// notifications.
	IntegrationIDs []string `json:"integration_ids"`
	// MappingStatus indicates whether the mapping is permitted to send
	// notifications. One of MappingStatusActive or MappingStatusInactive.
	MappingStatus string `json:"mapping_status"`
	// CreatedAt is when the mapping was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the mapping was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// Collection is the API response structure for a single page of the list
// operation.
type Collection struct {
	// Data is the list of returned event mappings.
	Data []EventMapping `json:"data"`
	// Meta contains pagination metadata.
	Meta Meta `json:"meta"`
}

// Meta is the pagination metadata returned by the list operation.
type Meta struct {
	// Total is the total number of results matching the filters.
	Total int `json:"total"`
	// Limit is the maximum number of results returned per page.
	Limit int `json:"limit"`
	// Sort is the sort order applied to the results.
	Sort string `json:"sort"`
	// NextCursor is the cursor used to retrieve the next page of results.
	// Nil when the results are the last page.
	NextCursor *string `json:"next_cursor,omitempty"`
	// PreviousCursor is the cursor used to retrieve the previous page of
	// results. Omitted when the results are the first page.
	PreviousCursor string `json:"previous_cursor,omitempty"`
}
