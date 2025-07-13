package ngwaf

import (
	"context"

	"github.com/fastly/go-fastly/v11/fastly"
	"github.com/fastly/go-fastly/v11/fastly/products"
	"github.com/fastly/go-fastly/v11/internal/productcore"
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

// EnableOutput holds the details returned by the API from 'Get' and
// 'Enable' operations; this alias exists to ensure that users of this
// package will have a stable name to reference.
type EnableOutput = products.EnableOutput

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

// ConfigureOutput holds the details returned by the API from
// 'GetConfiguration' and 'UpdateConfiguration' operations.
type ConfigureOutput struct {
	products.ConfigureOutput `mapstructure:",squash"`
	Configuration            *configureOutputNested `mapstructure:"configuration"`
}

type configureOutputNested struct {
	WorkspaceID *string `mapstructure:"workspace_id,omitempty"`
	TrafficRamp *string `mapstructure:"traffic_ramp,omitempty"`
}

// Get gets the status of the Next-Gen WAF product on the service.
func Get(ctx context.Context, c *fastly.Client, serviceID string) (EnableOutput, error) {
	return productcore.Get[EnableOutput](ctx, &productcore.GetInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Enable enables the Next-Gen WAF product on the service.
func Enable(ctx context.Context, c *fastly.Client, serviceID string, i EnableInput) (EnableOutput, error) {
	if i.WorkspaceID == "" {
		return EnableOutput{}, ErrMissingWorkspaceID
	}

	return productcore.Put[EnableOutput](ctx, &productcore.PutInput[EnableInput]{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
		Input:     i,
	})
}

// Disable disables the Next-Gen WAF product on the service.
func Disable(ctx context.Context, c *fastly.Client, serviceID string) error {
	return productcore.Delete(ctx, &productcore.DeleteInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// GetConfiguration gets the configuration of the Next-Gen WAF product on the service.
func GetConfiguration(ctx context.Context, c *fastly.Client, serviceID string) (ConfigureOutput, error) {
	return productcore.Get[ConfigureOutput](ctx, &productcore.GetInput{
		Client:        c,
		ProductID:     ProductID,
		ServiceID:     serviceID,
		URLComponents: []string{"configuration"},
	})
}

// UpdateConfiguration updates the configuration of the Next-Gen WAF product on the service.
func UpdateConfiguration(ctx context.Context, c *fastly.Client, serviceID string, i ConfigureInput) (ConfigureOutput, error) {
	return productcore.Patch[ConfigureOutput](ctx, &productcore.PatchInput[ConfigureInput]{
		Client:        c,
		ProductID:     ProductID,
		ServiceID:     serviceID,
		URLComponents: []string{"configuration"},
		Input:         i,
	})
}
