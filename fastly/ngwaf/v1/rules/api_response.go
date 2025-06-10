package rules

import (
	"encoding/json"
	"errors"
	"time"
)

// Scope defines the scope of a rule, including its type and the workspaces it applies to.
type Scope struct {
	// Type indicates the scope type (e.g., "workspace").
	Type string `json:"type"`
	// AppliesTo lists the workspace IDs that the rule applies to.
	AppliesTo []string `json:"applies_to"`
}

// ConditionItem represents a condition used in a rule, which can be either single or grouped.
type ConditionItem struct {
	// Type indicates whether the condition is "single" or "group".
	Type string `json:"type"`
	// Fields contains the condition-specific details, either SingleCondition or GroupCondition.
	Fields any
}

// SingleCondition defines a basic condition based on a field, operator, and value.
type SingleCondition struct {
	// Field is the request attribute to evaluate (e.g., "ip", "path").
	Field string `json:"field"`
	// Operator is the comparison operator (e.g., "equals", "contains").
	Operator string `json:"operator"`
	// Value is the comparison value for the field.
	Value string `json:"value"`
}

// GroupCondition defines a set of conditions and how they are logically grouped.
type GroupCondition struct {
	// GroupOperator specifies how to evaluate the conditions (e.g., "any", "all").
	GroupOperator string `json:"group_operator"`
	// Conditions lists the nested single conditions within this group.
	Conditions []Condition `json:"conditions"`
}

// Condition is a simplified form used inside a group condition.
type Condition struct {
	// Type specifies the condition type (should be "single").
	Type string `json:"type"`
	// Field is the request attribute to evaluate.
	Field string `json:"field"`
	// Operator is the comparison operator.
	Operator string `json:"operator"`
	// Value is the value to compare against.
	Value string `json:"value"`
}

// Action defines the action that will be executed when a rule is triggered.
type Action struct {
	// Type specifies the action type (e.g., "redirect", "block").
	Type string `json:"type"`
	// Signal is used when the action type is "exclude_signal".
	Signal string `json:"signal"`
	// RedirectURL is the URL to redirect to when using a redirect action.
	RedirectURL string `json:"redirect_url"`
	// ResponseCode is the HTTP response code to use for the redirect.
	ResponseCode int `json:"response_code"`
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
	// GroupOperator is the top-level logical operator (e.g., "any", "all").
	GroupOperator string `json:"group_operator"`
	// RequestLogging controls if requests matching this rule are logged ("none" or "sampled").
	RequestLogging string `json:"request_logging"`
	// Conditions defines the main conditions used to match requests.
	Conditions []ConditionItem `json:"conditions"`
	// Actions lists the actions to take when the rule is triggered.
	Actions []Action `json:"actions"`
	// CreatedAt is the timestamp when the rule was created.
	CreatedAt time.Time `mapstructure:"created_at" json:"created_at"`
	// UpdatedAt is the timestamp when the rule was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// UnmarshalJSON handles deserialization of ConditionItem, distinguishing between single and group conditions.
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
