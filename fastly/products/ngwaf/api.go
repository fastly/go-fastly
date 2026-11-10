package ngwaf

import (
	"context"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/internal/serviceproductcore"
)

const (
	ProductID   = "ngwaf"
	ProductName = "Next-Gen WAF"
)

// EnableInput holds the details required by the API's 'Enable'
// operation.
type EnableInput struct {
	WorkspaceID string `json:"workspace_id"`
}

// ServiceEnableOutput holds the details returned by the API from 'Get' and
// 'Enable' operations; this alias exists to ensure that users of this
// package will have a stable name to reference.
type ServiceEnableOutput = products.ServiceEnableOutput

// ErrMissingWorkspaceID is the error returned by the Enable function
// when it is passed an EnableInput struct with a WorkspaceID field
// that is empty.
var ErrMissingWorkspaceID = fastly.NewFieldError("WorkspaceID")

// ConfigureInput holds the details required by the API's
// 'UpdateConfiguration' operation.
type ConfigureInput struct {
	WorkspaceID string `json:"workspace_id,omitempty"`
	TrafficRamp string `json:"traffic_ramp,omitempty"`
}

// ServiceConfigureOutput holds the details returned by the API from
// 'GetConfiguration' and 'UpdateConfiguration' operations.
type ServiceConfigureOutput struct {
	products.ServiceConfigureOutput `mapstructure:",squash"`
	Configuration                   *configureOutputNested `mapstructure:"configuration"`
}

type configureOutputNested struct {
	WorkspaceID *string `mapstructure:"workspace_id,omitempty"`
	TrafficRamp *string `mapstructure:"traffic_ramp,omitempty"`
}

// Get gets the status of the Next-Gen WAF product on the service.
func Get(ctx context.Context, c *fastly.Client, serviceID string) (ServiceEnableOutput, error) {
	return serviceproductcore.Get[ServiceEnableOutput](ctx, &serviceproductcore.GetInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Enable enables the Next-Gen WAF product on the service.
func Enable(ctx context.Context, c *fastly.Client, serviceID string, i EnableInput) (ServiceEnableOutput, error) {
	if i.WorkspaceID == "" {
		return ServiceEnableOutput{}, ErrMissingWorkspaceID
	}

	return serviceproductcore.Put[ServiceEnableOutput](ctx, &serviceproductcore.PutInput[EnableInput]{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
		Input:     i,
	})
}

// Disable disables the Next-Gen WAF product on the service.
func Disable(ctx context.Context, c *fastly.Client, serviceID string) error {
	return serviceproductcore.Delete(ctx, &serviceproductcore.DeleteInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// GetConfiguration gets the configuration of the Next-Gen WAF product on the service.
func GetConfiguration(ctx context.Context, c *fastly.Client, serviceID string) (ServiceConfigureOutput, error) {
	return serviceproductcore.Get[ServiceConfigureOutput](ctx, &serviceproductcore.GetInput{
		Client:        c,
		ProductID:     ProductID,
		ServiceID:     serviceID,
		URLComponents: []string{"configuration"},
	})
}

// UpdateConfiguration updates the configuration of the Next-Gen WAF product on the service.
func UpdateConfiguration(ctx context.Context, c *fastly.Client, serviceID string, i ConfigureInput) (ServiceConfigureOutput, error) {
	return serviceproductcore.Patch[ServiceConfigureOutput](ctx, &serviceproductcore.PatchInput[ConfigureInput]{
		Client:        c,
		ProductID:     ProductID,
		ServiceID:     serviceID,
		URLComponents: []string{"configuration"},
		Input:         i,
	})
}
