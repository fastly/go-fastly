package rules

import (
	"encoding/json"
	"errors"
	"time"
)

// Scope defines the scope of a rule, including its type and the
// workspaces it applies to.
type Scope struct {
	// Type indicates the scope type (e.g., "workspace").
	Type string `json:"type"`
	// AppliesTo lists the workspace IDs that the rule applies to.
	AppliesTo []string `json:"applies_to"`
}

// ConditionItem represents a condition used in a rule, which can be
// either single or grouped.
type ConditionItem struct {
	// Type indicates whether the condition is "single" or "group".
	Type string `json:"type"`
	// Fields contains the condition-specific details, either
	// SingleCondition or GroupCondition.
	Fields any
}

// SingleCondition defines a basic condition based on a field,
// operator, and value.
type SingleCondition struct {
	// Field is the request attribute to evaluate (e.g., "ip",
	// "path").
	Field string `json:"field"`
	// Operator is the comparison operator (e.g., "equals",
	// "contains").
	Operator string `json:"operator"`
	// Value is the comparison value for the field.
	Value string `json:"value"`
}

// GroupCondition defines a set of conditions and how they are
// logically grouped.
type GroupCondition struct {
	// GroupOperator specifies how to evaluate the conditions
	// (e.g., `any`, `all`).
	GroupOperator string `json:"group_operator"`
	// Conditions lists the nested single conditions within this
	// group.
	Conditions []Condition `json:"conditions"`
}

// MultivalCondition defines a set of conditions and how they are
// logically grouped.

type MultivalCondition struct {
	// Field is the request attribute to evaluate (e.g., "post_parameter",
	// "signal").
	Field string `json:"field"`
	// Operator is the comparison operator (e.g., "equals",
	// "contains").
	Operator string `json:"operator"`
	// GroupOperator specifies how to evaluate the conditions
	// (e.g., `any`, `all`).
	GroupOperator string `json:"group_operator"`
	// Conditions lists the nested single conditions within this
	// group.
	Conditions []ConditionMul `json:"conditions"`
}

// Condition is a simplified form used inside a group condition.
type Condition struct {
	// Type specifies the condition type (should be `single`).
	Type string `json:"type"`
	// Field is the request attribute to evaluate.
	Field string `json:"field"`
	// Operator is the comparison operator.
	Operator string `json:"operator"`
	// Value is the value to compare against.
	Value string `json:"value"`
}

// ConditionMul is a simplified form used inside a multival condition.
type ConditionMul struct {
	// Type specifies the condition type (should be `single`).
	Type string `json:"type"`
	// Field is the request attribute to evaluate.
	Field string `json:"field"`
	// Operator is the comparison operator.
	Operator string `json:"operator"`
	// Value is the value to compare against.
	Value string `json:"value"`
}

// Action defines the action that will be executed when a rule is
// triggered.
type Action struct {
	// AllowInteractive specifies if interaction is allowed and is
	// only used for the "browser_challenge" action
	AllowInteractive *bool `json:"allow_interactive,omitempty"`
	// DeceptionType specifies the type of deception and is only
	// used for the "deception" action
	DeceptionType string `json:"deception_type"`
	// RedirectURL is the URL to redirect to when using a redirect
	// action.
	RedirectURL string `json:"redirect_url"`
	// ResponseCode is the HTTP response code to use for the
	// redirect.
	ResponseCode int `json:"response_code"`
	// Signal is used when the action type is "exclude_signal".
	Signal string `json:"signal"`
	// Type specifies the action type (e.g., "redirect", "block").
	Type string `json:"type"`
}

// RateLimit is the parameters of the rate limit rule.
type RateLimit struct {
	// List of client identifiers used for rate limiting.
	ClientIdentifiers []ClientIdentifier `json:"client_identifiers"`
	// Duration in seconds for the rate limit.
	Duration int `json:"duration"`
	// Time interval for the rate limit in seconds (60, 600, or 3600 minutes).
	Interval int `json:"interval"`
	// The signal used to count requests.
	Signal string `json:"signal"`
	// Rate limit threshold (between 1 and 10000).
	Threshold int `json:"threshold"`
}

// ClientIdentifiers defines how a client is identified.
type ClientIdentifier struct {
	// Key is the of the client identifier
	Key string `json:"key"`
	// Name is the name of the client identifier
	Name string `json:"name"`
	// Type is the type of identifier.
	Type string `json:"type"`
}

// Rule represents the complete configuration of a WAF rule.
type Rule struct {
	// RuleID is the unique identifier of the rule.
	RuleID string `json:"id"`
	// Type is the rule type (e.g., "request").
	Type string `json:"type"`
	// Scope indicates where the rule is applied.
	Scope Scope `json:"scope"`
	// Enabled indicates whether the rule is currently active.
	Enabled bool `json:"enabled"`
	// Description is a human-readable explanation of the rule.
	Description string `json:"description"`
	// GroupOperator is the top-level logical operator (e.g.,
	// `any`, `all`).
	GroupOperator string `json:"group_operator"`
	// RequestLogging controls if requests matching this rule are
	// logged (`none` or `sampled`).
	RequestLogging string `json:"request_logging"`
	// Conditions defines the main conditions used to match
	// requests.
	Conditions []ConditionItem `json:"conditions"`
	// Actions lists the actions to take when the rule is
	// triggered.
	Actions []Action `json:"actions"`
	// RateLimit defines the rate limit described by the rule.
	RateLimit *RateLimit `json:"rate_limit"`
	// CreatedAt is the timestamp when the rule was created.
	CreatedAt time.Time `mapstructure:"created_at" json:"created_at"`
	// UpdatedAt is the timestamp when the rule was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// UnmarshalJSON handles deserialization of ConditionItem,
// distinguishing between single, group and multival conditions.
func (ci *ConditionItem) UnmarshalJSON(data []byte) error {
	type alias struct {
		Type string `json:"type"`
	}
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	ci.Type = a.Type

	switch a.Type {
	case "single":
		var sc SingleCondition
		if err := json.Unmarshal(data, &sc); err != nil {
			return err
		}
		ci.Fields = sc
	case "group":
		var gc GroupCondition
		if err := json.Unmarshal(data, &gc); err != nil {
			return err
		}
		ci.Fields = gc
	case "multival":
		var mc MultivalCondition
		if err := json.Unmarshal(data, &mc); err != nil {
			return err
		}
	default:
		return errors.New("unknown condition type: " + a.Type)
	}

	return nil
}

// Rules is the response returned when listing multiple rules.
type Rules struct {
	// Data is the list of rules.
	Data []Rule `json:"data"`
	// Meta contains pagination or summary metadata.
	Meta MetaRules `json:"meta"`
}

// MetaRules holds metadata about the rules list.
type MetaRules struct {
	// Limit is the maximum number of rules returned.
	Limit int `json:"limit"`
	// Total is the total number of rules available.
	Total int `json:"total"`
}
