package eventtypes

// EventType describes a supported audit event type.
type EventType struct {
	// EventType is the event type identifier. Use this value in the
	// EventTypes field when creating a mapping.
	EventType string `json:"event_type"`
	// DisplayName is a human-readable name for the event type.
	DisplayName string `json:"display_name"`
	// ScopeTypes is the scope types this event type is compatible with.
	ScopeTypes []string `json:"scope_types"`
}

// Collection is the API response structure for the list operation.
type Collection struct {
	// Data is the list of returned event types.
	Data []EventType `json:"data"`
}
