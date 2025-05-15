package workspaces

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// AttackSignalThresholdsUpdateInput are the parameters for system site alerts.
type AttackSignalThresholdsUpdateInput struct {
	OneMinute  *int  `json:"one_minute,omitempty"`
	TenMinutes *int  `json:"ten_minutes,omitempty"`
	OneHour    *int  `json:"one_hour,omitempty"`
	Immediate  *bool `json:"immediate,omitempty"`
}

// UpdateInput specifies the information needed for the Update() function to
// perform the operation.
type UpdateInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `json:"-"`
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string `json:"-"`
	// Name is the name of a workspace to create (required).
	Name *string `json:"name"`
	// Description is the description of a workspace.
	Description *string `json:"description"`
	// Mode is the mode of a workspace.
	Mode *string `json:"mode"`
	// AttackSignalThresholds are the parameters for system site alerts.
	AttackSignalThresholds *AttackSignalThresholdsUpdateInput `json:"attack_signal_thresholds"`
	// IPAnonymization is the selected option to anonymize IP addresses.
	IPAnonymization *string `json:"ip_anonymization"`
}

// Update updates the specified workspace.
func Update(c *fastly.Client, i *UpdateInput) (*Workspace, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID)

	resp, err := c.PatchJSON(path, i, fastly.CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ws *Workspace
	if err := json.NewDecoder(resp.Body).Decode(&ws); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}
	return ws, nil
}
