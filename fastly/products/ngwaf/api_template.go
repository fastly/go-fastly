//go:generate rm -f api.go api_test.go
//go:generate service_linked_product

package ngwaf

import (
	"github.com/fastly/go-fastly/v9/fastly"
)

const (
	ProductName = "Next-Gen WAF"
	ProductID   = "ngwaf"
)

// ErrMissingWorkspaceID is the error returned by the Enable function
// when it is passed an EnableInput struct with a WorkspaceID field
// that is empty.
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
