//go:generate service_linked_product

package ngwaf

import (
	"github.com/fastly/go-fastly/v9/fastly"
)

const (
	ProductName = "NextGenWAF"
	ProductID   = "ngwaf"
)

// ErrMissingWorkspaceID is an error that is returned when an input struct
// requires a "WorkspaceID" key, but one was not set.
var ErrMissingWorkspaceID = fastly.NewFieldError("WorkspaceID")

type EnableInput struct {
	WorkspaceID string `json:"workspace_id"`
}

func (i *EnableInput) Validate() error {
	if i.WorkspaceID == "" {
		return ErrMissingWorkspaceID
	}
	return nil
}

type ConfigureInput struct {
	WorkspaceID string `json:"workspace_id,omitempty"`
	TrafficRamp string `json:"traffic_ramp,omitempty"`
}

func (i *ConfigureInput) Validate() error {
	return nil
}

type ConfigureOutput struct {
	fastly.ProductConfiguration
	Configuration *ConfigureOutputNested `mapstructure:"configuration"`
}

type ConfigureOutputNested struct {
	WorkspaceID *string `mapstructure:"workspace_id,omitempty"`
	TrafficRamp *string `mapstructure:"traffic_ramp,omitempty"`
}
