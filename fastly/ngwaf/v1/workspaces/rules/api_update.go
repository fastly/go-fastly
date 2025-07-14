package rules

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fastly/go-fastly/v10/fastly"
)

// UpdateInput specifies the information needed for the Update()
// function to perform the operation.
type UpdateInput struct {
	// Type specifies the category of the rule (e.g., "request")
	// (required).
	Type *string
	// Description provides a human-readable explanation of what
	// the rule does (required).
	Description *string
	// GroupOperator defines the logical operator ("any" or "all")
	// used to evaluate grouped conditions.
	GroupOperator *string
	// Enabled determines if the rule is active. If false or
	// omitted, the rule is disabled by default.
	Enabled *bool
	// ExpiresAt sets a specific time when the rule will
	// automatically be disabled.
	ExpiresAt *time.Time
	// Actions is a list of actions that should be executed when
	// rule conditions are met (required).
	Actions []*UpdateAction
	// Conditions contains individual (non-grouped) matching
	// criteria.
	Conditions []*UpdateCondition
	// GroupConditions is a list of grouped conditions with nested
	// logical evaluation.
	GroupConditions []*UpdateGroupCondition
	// RequestLogging defines how request logs are handled when
	// the rule is matched ("sampled" or "none"). Applicable only
	// for request-type rules.
	RequestLogging *string
	// RuleID is the rule identifier (required).
	RuleID *string
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// UpdateAction represents an action taken when a rule's conditions
// are met.
type UpdateAction struct {
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

// Update updates a rule.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*Rule, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.RuleID == nil {
		return nil, fastly.ErrMissingRuleID
	}

	var mergedConditions []any
	for _, c := range i.Conditions {
		mergedConditions = append(mergedConditions, c)
	}
	for _, gc := range i.GroupConditions {
		mergedConditions = append(mergedConditions, gc)
	}

	v := struct {
		Type           *string         `json:"type,omitempty"`
		Description    *string         `json:"description,omitempty"`
		GroupOperator  *string         `json:"group_operator,omitempty"`
		Enabled        *bool           `json:"enabled,omitempty"`
		ExpiresAt      *time.Time      `json:"expires_at,omitempty"`
		Actions        []*UpdateAction `json:"actions,omitempty"`
		Conditions     []any           `json:"conditions,omitempty"`
		RequestLogging *string         `json:"request_logging,omitempty"`
	}{
		Type:           i.Type,
		Description:    i.Description,
		GroupOperator:  i.GroupOperator,
		Enabled:        i.Enabled,
		ExpiresAt:      i.ExpiresAt,
		Actions:        i.Actions,
		Conditions:     mergedConditions,
		RequestLogging: i.RequestLogging,
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "rules", *i.RuleID)

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
