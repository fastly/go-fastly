package scope

// Type defines the type of scope.
type Type string

const (
	ScopeTypeWorkspace Type = "workspace"
	ScopeTypeAccount   Type = "account"
)

// Scope defines the scope of a resource, specifying its type and applicable workspace identifiers.
//
// The `Type` field determines whether the resource is scoped to a "workspace" or "account":
//   - "workspace": The rule applies only to a specific workspace.
//   - "account": The rule applies at the account level and may span multiple workspaces.
//
// The `AppliesTo` field is an array of workspace IDs:
//   - For "workspace" scope, this must contain exactly one workspace ID.
//   - For "account" scope, this can include multiple workspace IDs or a wildcard "*" to target all workspaces.
type Scope struct {
	// Type specifies whether the rule applies at the "workspace" or "account" level (required).
	Type Type `json:"type"`
	// AppliesTo lists the workspace IDs the rule applies to. For "workspace" type, a single ID is required.
	// For "account" type, multiple IDs or "*" (wildcard for all workspaces) can be provided.
	AppliesTo []string `json:"applies_to"`
}
