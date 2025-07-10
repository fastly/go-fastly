package workspaces

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// AttackSignalThresholdsCreateInput are the parameters for system site alerts.
type AttackSignalThresholdsCreateInput struct {
	OneMinute  *int  `json:"one_minute,omitempty"`
	TenMinutes *int  `json:"ten_minutes,omitempty"`
	OneHour    *int  `json:"one_hour,omitempty"`
	Immediate  *bool `json:"immediate,omitempty"`
}

// CreateInput specifies the information needed for the Create() function to
// perform the operation.
type CreateInput struct {
	// AttackSignalThresholds are the parameters for system site alerts.
	AttackSignalThresholds *AttackSignalThresholdsCreateInput `json:"attack_signal_thresholds,omitempty"`
	// ClientIPHeaders lists the request headers containing the client IP address.
	ClientIPHeaders []string `json:"client_ip_headers,omitempty"`
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `json:"-"`
	// DefaultBlockingResponseCode is the default response code.
	DefaultBlockingResponseCode *int `json:"default_blocking_response_code,omitempty"`
	// Description is the description of a workspace.
	Description *string `json:"description,omitempty"`
	// IPAnonymization is the selected option to anonymize IP addresses.
	IPAnonymization *string `json:"ip_anonymization,omitempty"`
	// Mode is the mode of a workspace (required).
	Mode *string `json:"mode"`
	// Name is the name of a workspace to create (required).
	Name *string `json:"name"`
}

// Create creates a new workspace.
func Create(c *fastly.Client, i *CreateInput) (*Workspace, error) {
	if i.Name == nil {
		return nil, fastly.ErrMissingName
	}

	if i.Mode == nil {
		return nil, fastly.ErrMissingMode
	}

	resp, err := c.PostJSON("/ngwaf/v1/workspaces", i, fastly.CreateRequestOptions(i.Context))
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
