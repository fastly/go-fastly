package common

import (
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// BuildPath generates the appropriate API path based on the given scope,
// resource type, and resource ID.
//
// For scope type "workspace", the first element of AppliesTo is used as the workspace ID.
// For scope type "account", the workspace ID is ignored.
func BuildPath(scope *Scope, resource, resourceID string) (string, error) {
	if scope == nil {
		return "", fmt.Errorf("scope is required to build path")
	}

	switch scope.Type {
	case ScopeTypeWorkspace:
		if len(scope.AppliesTo) > 1 || scope.AppliesTo[0] == "" {
			return "", fmt.Errorf("scope.applies_to must contain exactly one workspace ID for scope type 'workspace'")
		}
		workspaceID := scope.AppliesTo[0]

		if resourceID == "" {
			return fastly.ToSafeURL("ngwaf", "v1", "workspaces", workspaceID, resource), nil
		}
		return fastly.ToSafeURL("ngwaf", "v1", "workspaces", workspaceID, resource, resourceID), nil

	case ScopeTypeAccount:
		if resourceID == "" {
			return fastly.ToSafeURL("ngwaf", "v1", resource), nil
		}
		return fastly.ToSafeURL("ngwaf", "v1", resource, resourceID), nil

	default:
		return "", fmt.Errorf("unsupported scope type: %s", scope.Type)
	}
}
