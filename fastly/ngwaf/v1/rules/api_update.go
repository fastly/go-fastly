package rules

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fastly/go-fastly/v11/fastly"
	"github.com/fastly/go-fastly/v11/fastly/ngwaf/v1/common"
)

// UpdateInput specifies the information needed for the Update()
// function to perform the operation.
type UpdateInput struct {
	// Actions is a list of actions that should be executed when
	// rule conditions are met (required).
	Actions []*UpdateAction
	// Conditions contains individual (non-grouped) matching
	// criteria.
	Conditions []*UpdateCondition
	// Description provides a human-readable explanation of what
	// the rule does (required).
	Description *string
	// Enabled determines if the rule is active. If false or
	// omitted, the rule is disabled by default.
	Enabled *bool
	// ExpiresAt sets a specific time when the rule will
	// automatically be disabled.
	ExpiresAt *time.Time
	// GroupConditions is a list of grouped conditions with nested
	// logical evaluation.
	GroupConditions []*UpdateGroupCondition
	// GroupOperator defines the logical operator ("any" or "all")
	// used to evaluate grouped conditions.
	GroupOperator *string
	// RateLimit defines how rate limit rules are enforced.
	RateLimit *UpdateRateLimit
	// RequestLogging defines how request logs are handled when
	// the rule is matched ("sampled" or "none"). Applicable only
	// for request-type rules.
	RequestLogging *string
	// RuleID is the rule identifier (required).
	RuleID *string
	// Scope defines where the rule is applied, including its type
	// (e.g., "workspace" or "account") and the specific IDs it
	// applies to (required).
	Scope *common.Scope
	// Type specifies the category of the rule (e.g., "request")
	// (required).
	Type *string
}

// UpdateAction represents an action taken when a rule's conditions
// are met.
type UpdateAction struct {
	// AllowInteractive specifies if interaction is allowed and is
	// only used for the "browser_challenge" action
	AllowInteractive *bool `json:"allow_interactive,omitempty"`
	// DeceptionType specifies the type of deception and is only
	// used for the "deception" action
	DeceptionType *string `json:"deception_type,omitempty"`
	// RedirectURL specifies the target URL when redirecting the
	// request.
	RedirectURL *string `json:"redirect_url,omitempty"`
	// ResponseCode is the HTTP status code returned during
	// redirection (e.g., 301, 302).
	ResponseCode *int `json:"response_code,omitempty"`
	// Signal is the signal name used only for the
	// "exclude_signal" action type.
	Signal *string `json:"signal,omitempty"`
	// Type specifies the action type (e.g., "block",
	// "exclude_signal") (required).
	Type *string `json:"type"`
}

// UpdateCondition defines a single condition.
type UpdateCondition struct {
	// Type specifies the condition type (must be "single")
	// (required).
	Type *string `json:"type"`
	// Field is the name of the field to be evaluated (e.g., "ip",
	// "path") (required).
	Field *string `json:"field"`
	// Operator determines how the field is evaluated (e.g.,
	// "equals", "contains") (required).
	Operator *string `json:"operator"`
	// Value is the value against which the field is compared
	// (required).
	Value *string `json:"value"`
}

// UpdateGroupCondition defines a group of conditions with a logical
// operator.
type UpdateGroupCondition struct {
	// Type specifies the condition group type (must be "group")
	// (required).
	Type *string `json:"type"`
	// GroupOperator is the logical operator used to evaluate the
	// conditions ("any" or "all") (required).
	GroupOperator *string `json:"group_operator"`
	// Conditions is the list of single conditions to evaluate
	// within the group (required).
	Conditions []*UpdateCondition `json:"conditions"`
}

// UpdateRateLimit defines how rate limit rules are enforced.
type UpdateRateLimit struct {
	// List of client identifiers used for rate limiting. Can only be length 1 or 2.
	ClientIdentifiers []*UpdateClientIdentifier
	// Duration in seconds for the rate limit.
	Duration *int
	// Time interval for the rate limit in seconds (60, 600, or 3600 minutes).
	Interval *int
	// The signal used to count requests.
	Signal *string
	// Rate limit threshold (between 1 and 10000).
	Threshold *int
}

// UpdateClientIdentifier is the client identifier for rate limit rules.
type UpdateClientIdentifier struct {
	// Key is the of the client identifier
	Key *string `json:"key,omitempty"`
	// Name is the name of the client identifier
	Name *string `json:"name,omitempty"`
	// Type is the type of the client identifier
	Type *string `json:"type"`
}

// Update updates a rule.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*Rule, error) {
	if i.RuleID == nil {
		return nil, fastly.ErrMissingRuleID
	}
	if i.Scope == nil {
		return nil, fastly.ErrMissingScope
	}

	var mergedConditions []any
	for _, c := range i.Conditions {
		mergedConditions = append(mergedConditions, c)
	}
	for _, gc := range i.GroupConditions {
		mergedConditions = append(mergedConditions, gc)
	}

	v := struct {
		Actions        []*UpdateAction  `json:"actions,omitempty"`
		Conditions     []any            `json:"conditions,omitempty"`
		Description    *string          `json:"description,omitempty"`
		Enabled        *bool            `json:"enabled,omitempty"`
		ExpiresAt      *time.Time       `json:"expires_at,omitempty"`
		GroupOperator  *string          `json:"group_operator,omitempty"`
		RateLimit      *UpdateRateLimit `json:"rate_limit,omitempty"`
		RequestLogging *string          `json:"request_logging,omitempty"`
		Scope          *common.Scope    `json:"scope,omitempty"`
		Type           *string          `json:"type,omitempty"`
	}{
		Actions:        i.Actions,
		Conditions:     mergedConditions,
		Description:    i.Description,
		Enabled:        i.Enabled,
		ExpiresAt:      i.ExpiresAt,
		GroupOperator:  i.GroupOperator,
		RequestLogging: i.RequestLogging,
		RateLimit:      i.RateLimit,
		Scope:          i.Scope,
		Type:           i.Type,
	}

	path, err := common.BuildPath(i.Scope, "rules", *i.RuleID)
	if err != nil {
		return nil, fmt.Errorf("failed to build API path: %w", err)
	}

	resp, err := c.PatchJSON(ctx, path, v, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r *Rule
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return r, nil
}
