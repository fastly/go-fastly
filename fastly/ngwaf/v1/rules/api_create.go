package rules

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fastly/go-fastly/v10/fastly"
	"github.com/fastly/go-fastly/v10/fastly/ngwaf/v1/common"
)

// CreateInput specifies the information needed for the Create()
// function to perform the operation.
type CreateInput struct {
	// Actions is a list of actions that should be executed when rule conditions are met (required).
	Actions []*CreateAction
	// Conditions contains individual (non-grouped) matching
	// criteria.
	Conditions []*CreateCondition
	// Description provides a human-readable explanation of what the rule does (required).
	Description *string
	// Enabled determines if the rule is active. If false or omitted, the rule is disabled by default.
	Enabled *bool
	// ExpiresAt sets a specific time when the rule will automatically be disabled.
	ExpiresAt *time.Time
	// GroupConditions is a list of grouped conditions with nested logical evaluation.
	GroupConditions []*CreateGroupCondition
	// GroupOperator defines the logical operator ("any" or "all") used to evaluate grouped conditions.
	GroupOperator *string
	// RequestLogging defines how request logs are handled when the rule is matched ("sampled" or "none"). Applicable only for request-type rules.
	RequestLogging *string
	// Scope defines where the rule is applied, including its type (e.g., "workspace" or "account") and the specific IDs it applies to (required).
	Scope *common.Scope
	// Type specifies the category of the rule (e.g., "request") (required).
	Type *string
}

// CreateAction represents an action taken when a rule's conditions
// are met.
type CreateAction struct {
	// Type specifies the action type (e.g., "block",
	// "exclude_signal") (required).
	Type *string `json:"type"`
	// Signal is the signal name used only for the
	// "exclude_signal" action type.
	Signal *string `json:"signal,omitempty"`
	// RedirectURL specifies the target URL when redirecting the
	// request.
	RedirectURL *string `json:"redirect_url,omitempty"`
	// ResponseCode is the HTTP status code returned during
	// redirection (e.g., 301, 302).
	ResponseCode *int `json:"response_code,omitempty"`
}

// CreateCondition defines a single condition.
type CreateCondition struct {
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

// CreateGroupCondition defines a group of conditions with a logical
// operator.
type CreateGroupCondition struct {
	// Type specifies the condition group type (must be "group")
	// (required).
	Type *string `json:"type"`
	// GroupOperator is the logical operator used to evaluate the
	// conditions ("any" or "all") (required).
	GroupOperator *string `json:"group_operator"`
	// Conditions is the list of single conditions to evaluate
	// within the group (required).
	Conditions []*CreateCondition `json:"conditions"`
}

// Create creates a new rule.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*Rule, error) {
	if i.Type == nil {
		return nil, fastly.ErrMissingType
	}
	if i.Description == nil {
		return nil, fastly.ErrMissingDescription
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
	if len(mergedConditions) == 0 {
		return nil, fastly.ErrMissingConditions
	}

	v := struct {
		Actions        []*CreateAction `json:"actions"`
		Conditions     []any           `json:"conditions"`
		Description    *string         `json:"description"`
		Enabled        *bool           `json:"enabled,omitempty"`
		ExpiresAt      *time.Time      `json:"expires_at,omitempty"`
		GroupOperator  *string         `json:"group_operator,omitempty"`
		RequestLogging *string         `json:"request_logging,omitempty"`
		Scope          *common.Scope   `json:"scope"`
		Type           *string         `json:"type"`
	}{
		Actions:        i.Actions,
		Conditions:     mergedConditions,
		Description:    i.Description,
		Enabled:        i.Enabled,
		ExpiresAt:      i.ExpiresAt,
		GroupOperator:  i.GroupOperator,
		RequestLogging: i.RequestLogging,
		Scope:          i.Scope,
		Type:           i.Type,
	}

	path, err := common.BuildPath(i.Scope, "rules", "")
	if err != nil {
		return nil, fmt.Errorf("failed to build API path: %w", err)
	}

	resp, err := c.PostJSON(ctx, path, v, fastly.CreateRequestOptions())
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
