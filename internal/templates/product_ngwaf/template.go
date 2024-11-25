//go:generate service_product_enablement

package product_ngwaf

import (
	fsterr "github.com/fastly/go-fastly/v9/pkg/errors"
)

const (
	ProductName = "NextGenWAF"
	ProductID   = "ngwaf"
)

type EnableInput struct {
	WorkspaceID string `json:"workspace_id"`
}

type ConfigureInput struct {
	WorkspaceID string `json:"workspace_id,omitempty"`
	TrafficRamp string `json:"traffic_ramp,omitempty"`
}

func (i *EnableInput) Validate() error {
	if i.WorkspaceID == "" {
		return fsterr.ErrMissingWorkspaceID
	}
	return nil
}

func (i *ConfigureInput) Validate() error {
	return nil
}
