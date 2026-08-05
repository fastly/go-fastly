package scopetypes

// ScopeType describes a supported mapping scope type.
type ScopeType struct {
	// ScopeType is the scope type identifier.
	ScopeType string `json:"scope_type"`
}

// Collection is the API response structure for the list operation.
type Collection struct {
	// Data is the list of returned scope types.
	Data []ScopeType `json:"data"`
}
