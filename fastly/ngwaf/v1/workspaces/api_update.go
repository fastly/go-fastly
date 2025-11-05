package workspaces

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v12/fastly"
)

// AttackSignalThresholdsUpdateInput are the parameters for system
// site alerts.
type AttackSignalThresholdsUpdateInput struct {
	OneMinute  *int  `json:"one_minute,omitempty"`
	TenMinutes *int  `json:"ten_minutes,omitempty"`
	OneHour    *int  `json:"one_hour,omitempty"`
	Immediate  *bool `json:"immediate,omitempty"`
}

// UpdateInput specifies the information needed for the Update()
// function to perform the operation.
type UpdateInput struct {
	// AttackSignalThresholds are the parameters for system site
	// alerts.
	AttackSignalThresholds *AttackSignalThresholdsUpdateInput `json:"attack_signal_thresholds,omitempty"`
	// ClientIPHeaders lists the request headers containing the
	// client IP address.
	ClientIPHeaders []string `json:"client_ip_headers,omitempty"`
	// DefaultBlockingResponseCode is the default response code.
	DefaultBlockingResponseCode *int `json:"default_blocking_response_code,omitempty"`
	// DefaultRedirectURL is a URL that will be used if
	// the DefaultBlockingResponseCode is set to 301 or 302
	DefaultRedirectURL *string `json:"default_redirect_url,omitempty"`
	// Description is the description of the workspace.
	Description *string `json:"description,omitempty"`
	// IPAnonymization is the selected option to anonymize IP
	// addresses.
	IPAnonymization *string `json:"ip_anonymization,omitempty"`
	// Mode is the mode of the workspace.
	Mode *string `json:"mode,omitempty"`
	// Name is the new name of the workspace.
	Name *string `json:"name,omitempty"`
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string `json:"-"`
}

// Update updates the specified workspace.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*Workspace, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID)

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
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
