package rules

import (
	"time"
)

// Collection is the API response structure for the list operation.
type Collection struct {
	// Data contains the API data.
	Data []Data `json:"data"`
	// Meta contains metadata related to paginating the full dataset.
	Meta Meta `json:"meta"`
}

// Data is a subset of the API response structure containing the specific API
// data itself.
type Data struct {
	// Action determines where matching requests are routed.
	Action Action `json:"action"`
	// Conditions are the matching criteria evaluated against the incoming
	// request. A rule with no conditions is the default (catch-all) rule for
	// its path.
	Conditions []Condition `json:"conditions"`
	// CreatedAt is the date and time in ISO 8601 format.
	CreatedAt time.Time `json:"created_at"`
	// IsDefault indicates whether this is the catch-all rule for its path
	// (i.e. it has no conditions). This is a read-only, computed field.
	IsDefault bool `json:"is_default"`
	// RuleID is the rule identifier.
	RuleID string `json:"id"`
	// UpdatedAt is the date and time in ISO 8601 format.
	UpdatedAt time.Time `json:"updated_at"`
}

// Action determines where a rule's matching requests are routed.
type Action struct {
	// Type is the action type (e.g. "service").
	Type string `json:"type"`
	// Value is the destination for the action (e.g. a service ID when Type is
	// "service").
	Value string `json:"value"`
}

// Condition is a single matching criterion evaluated against the incoming
// request.
type Condition struct {
	// Key is the header name to match against. Only applicable when Type is
	// "header".
	Key *string `json:"key,omitempty"`
	// Operator is the comparison operator used to evaluate Value against the
	// request (e.g. "equals").
	Operator string `json:"operator"`
	// Type is the condition category (e.g. "header").
	Type string `json:"type"`
	// Value is the value compared against the request using Operator.
	Value string `json:"value"`
}

// Meta is a subset of the API response structure containing metadata related to
// paginating the full dataset.
type Meta struct {
	// Limit is how many results are included in this response.
	Limit int `json:"limit"`
	// NextCursor is the cursor value used to retrieve the next page.
	NextCursor string `json:"next_cursor"`
	// Sort is the field used to order the response by.
	Sort string `json:"sort"`
	// Total is the total number of results.
	Total int `json:"total"`
}
