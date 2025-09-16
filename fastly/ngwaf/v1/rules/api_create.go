package rules

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fastly/go-fastly/v11/fastly"
	"github.com/fastly/go-fastly/v11/fastly/ngwaf/v1/scope"
)

// CreateInput specifies the information needed for the Create()
// function to perform the operation.
type CreateInput struct {
	// Actions is a list of actions that should be executed when
	// rule conditions are met (required).
	Actions []*CreateAction
	// Conditions contains individual (non-grouped) matching
	// criteria.
	Conditions []*CreateCondition
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
	GroupConditions []*CreateGroupCondition
	// MultivalConditions lists the nested multival conditions within this
	// group.
	MultivalConditions []*CreateMultivalCondition
	// GroupOperator defines the logical operator ("any" or "all")
	// used to evaluate grouped conditions.
	GroupOperator *string
	// RateLimit defines how rate limit rules are enforced.
	RateLimit *CreateRateLimit
	// RequestLogging defines how request logs are handled when
	// the rule is matched ("sampled" or "none"). Applicable only
	// for request-type rules.
	RequestLogging *string
	// Scope defines where the rule is applied, including its type
	// (e.g., "workspace" or "account") and the specific IDs it
	// applies to (required).
	Scope *scope.Scope
	// Type specifies the category of the rule (e.g., "request")
	// (required).
	Type *string
}

// CreateAction represents an action taken when a rule's conditions
// are met.
type CreateAction struct {
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

// CreateCondition defines a single condition.
type CreateCondition struct {
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

// CreateConditionMult defines a multival condition.
type CreateConditionMult struct {
	// Field is the name of the field to be evaluated (e.g., "name",
	// "value", "value_int") (required).
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
	// GroupOperator is the logical operator used to evaluate the
	// conditions ("any" or "all") (required).
	GroupOperator *string `json:"group_operator"`
	// Conditions is the list of single conditions to evaluate
	// within the group (required).
	Conditions []*CreateCondition `json:"conditions"`
}

// CreateMultivalCondition defines a multival of conditions with a logical
// operator.
type CreateMultivalCondition struct {
	// Field is the request attribute to evaluate (e.g., "post_parameter",
	// "signal").
	Field *string `json:"field"`
	// Operator is the comparison operator (e.g., "exists",
	// "does_not_exist").
	Operator *string `json:"operator"`
	// GroupOperator specifies how to evaluate the conditions
	// (e.g., `any`, `all`).
	GroupOperator *string `json:"group_operator"`
	// Conditions lists the nested single conditions within this
	// group.
	Conditions []*CreateConditionMult `json:"conditions"`
}

// CreateRateLimit defines how rate limit rules are enforced.
type CreateRateLimit struct {
	// List of client identifiers used for rate limiting. Can only be length 1 or 2.
	ClientIdentifiers []*CreateClientIdentifier `json:"client_identifiers"`
	// Duration in seconds for the rate limit.
	Duration *int `json:"duration"`
	// Time interval for the rate limit in seconds (60, 600, or 3600 minutes).
	Interval *int `json:"interval"`
	// The signal used to count requests.
	Signal *string `json:"signal"`
	// Rate limit threshold (between 1 and 10000).
	Threshold *int `json:"threshold"`
}

// CreateClientIdentifier is the client identifier for rate limit rules.
type CreateClientIdentifier struct {
	// Key is the of the client identifier
	Key *string `json:"key,omitempty"`
	// Name is the name of the client identifier
	Name *string `json:"name,omitempty"`
	// Type is the type of the client identifier
	Type *string `json:"type"`
}

// Private structs to ensure correct condition types.
type privateCreateCondition struct {
	Type     *string `json:"type"`
	Field    *string `json:"field"`
	Operator *string `json:"operator"`
	Value    *string `json:"value"`
}

type privateCreateConditionMult struct {
	Type     *string `json:"type"`
	Field    *string `json:"field"`
	Operator *string `json:"operator"`
	Value    *string `json:"value"`
}

type privateCreateGroupCondition struct {
	Type          *string                   `json:"type"`
	GroupOperator *string                   `json:"group_operator"`
	Conditions    []*privateCreateCondition `json:"conditions"`
}

type privateCreateMultivalCondition struct {
	Type          *string                       `json:"type"`
	Field         *string                       `json:"field"`
	Operator      *string                       `json:"operator"`
	GroupOperator *string                       `json:"group_operator"`
	Conditions    []*privateCreateConditionMult `json:"conditions"`
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
		privateCondition := &privateCreateCondition{
			Type:     fastly.ToPointer("single"),
			Field:    c.Field,
			Operator: c.Operator,
			Value:    c.Value,
		}
		mergedConditions = append(mergedConditions, privateCondition)
	}
	for _, gc := range i.GroupConditions {
		var privateSubConditions []*privateCreateCondition
		for _, subCond := range gc.Conditions {
			privateSubConditions = append(privateSubConditions, &privateCreateCondition{
				Type:     fastly.ToPointer("single"),
				Field:    subCond.Field,
				Operator: subCond.Operator,
				Value:    subCond.Value,
			})
		}
		privateGroupCondition := &privateCreateGroupCondition{
			Type:          fastly.ToPointer("group"),
			GroupOperator: gc.GroupOperator,
			Conditions:    privateSubConditions,
		}
		mergedConditions = append(mergedConditions, privateGroupCondition)
	}
	for _, mc := range i.MultivalConditions {
		var privateSubConditions []*privateCreateConditionMult
		for _, subCond := range mc.Conditions {
			privateSubConditions = append(privateSubConditions, &privateCreateConditionMult{
				Type:     fastly.ToPointer("single"),
				Field:    subCond.Field,
				Operator: subCond.Operator,
				Value:    subCond.Value,
			})
		}
		privateMultivalCondition := &privateCreateMultivalCondition{
			Type:          fastly.ToPointer("multival"),
			Field:         mc.Field,
			Operator:      mc.Operator,
			GroupOperator: mc.GroupOperator,
			Conditions:    privateSubConditions,
		}
		mergedConditions = append(mergedConditions, privateMultivalCondition)
	}
	if len(mergedConditions) == 0 {
		return nil, fastly.ErrMissingConditions
	}

	v := struct {
		Actions        []*CreateAction  `json:"actions,omitempty"`
		Conditions     []any            `json:"conditions"`
		Description    *string          `json:"description"`
		Enabled        *bool            `json:"enabled,omitempty"`
		ExpiresAt      *time.Time       `json:"expires_at,omitempty"`
		GroupOperator  *string          `json:"group_operator,omitempty"`
		RateLimit      *CreateRateLimit `json:"rate_limit,omitempty"`
		RequestLogging *string          `json:"request_logging,omitempty"`
		Scope          *scope.Scope     `json:"scope"`
		Type           *string          `json:"type"`
	}{
		Actions:        i.Actions,
		Conditions:     mergedConditions,
		Description:    i.Description,
		Enabled:        i.Enabled,
		ExpiresAt:      i.ExpiresAt,
		GroupOperator:  i.GroupOperator,
		RateLimit:      i.RateLimit,
		RequestLogging: i.RequestLogging,
		Scope:          i.Scope,
		Type:           i.Type,
	}

	path, err := scope.BuildPath(i.Scope, "rules", "")
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
