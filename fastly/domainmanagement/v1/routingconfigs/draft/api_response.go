package draft

import (
	"time"

	"github.com/fastly/go-fastly/v17/fastly/domainmanagement/v1/routingconfigs/paths"
	"github.com/fastly/go-fastly/v17/fastly/domainmanagement/v1/routingconfigs/paths/rules"
)

// Data is the API response structure for the update (comment) operation.
type Data struct {
	// Comment is the descriptive note associated with the draft version.
	Comment string `json:"comment"`
	// CreatedAt is the date and time in ISO 8601 format.
	CreatedAt time.Time `json:"created_at"`
	// VersionID is the identifier of the draft version.
	VersionID string `json:"id"`
}

// Diff is the API response structure for the diff operation, describing the
// differences between a routing config's active and draft versions.
type Diff struct {
	// Added lists paths present in the draft but not the active version.
	Added []PathChange `json:"added"`
	// Deleted lists paths present in the active version but not the draft.
	Deleted []PathChange `json:"deleted"`
	// Modified lists paths that exist in both versions but differ.
	Modified []PathModification `json:"modified"`
}

// PathChange describes a path that was wholly added or deleted between the
// active and draft versions.
type PathChange struct {
	// Path is the added or deleted path.
	Path *paths.Data `json:"path,omitempty"`
	// Rules lists the rules belonging to the added or deleted path.
	Rules []rules.Data `json:"rules,omitempty"`
}

// PathModification describes the differences for a path that exists in both
// the active and draft versions.
type PathModification struct {
	// OldPath is the path's previous URL pattern, present when it was renamed.
	OldPath string `json:"old_path,omitempty"`
	// Path is the path's current URL pattern.
	Path string `json:"path,omitempty"`
	// PathID is the identifier of the modified path.
	PathID string `json:"path_id"`
	// RulesAdded lists rules present in the draft but not the active version
	// for this path.
	RulesAdded []rules.Data `json:"rules_added,omitempty"`
	// RulesChanged lists rules that exist in both versions but differ for
	// this path.
	RulesChanged []RuleChange `json:"rules_changed,omitempty"`
	// RulesDeleted lists rules present in the active version but not the
	// draft for this path.
	RulesDeleted []rules.Data `json:"rules_deleted,omitempty"`
}

// RuleChange describes the differences for a rule that exists in both the
// active and draft versions.
type RuleChange struct {
	// NewAction is the rule's action in the draft version.
	NewAction *rules.Action `json:"new_action,omitempty"`
	// NewConditions is the rule's conditions in the draft version.
	NewConditions []rules.Condition `json:"new_conditions,omitempty"`
	// OldAction is the rule's action in the active version.
	OldAction *rules.Action `json:"old_action,omitempty"`
	// OldConditions is the rule's conditions in the active version.
	OldConditions []rules.Condition `json:"old_conditions,omitempty"`
	// RuleID is the identifier of the modified rule.
	RuleID string `json:"rule_id"`
}
