package ddosprotection

import (
	"context"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/internal/serviceproductcore"
)

const (
	ProductID   = "ddos_protection"
	ProductName = "DDoS Protection"
)

// EnableInput holds the details required by the API's 'Enable'
// operation.
type EnableInput struct {
	Mode string `json:"mode,omitempty"`
}

// EnableOutput holds the details returned by the API from 'Get' and
// 'Enable' operations; this alias exists to ensure that users of this
// package will have a stable name to reference.
type EnableOutput = products.ServiceEnableOutput

// ErrMissingMode is the error returned by the UpdateConfiguration
// function when it is passed a ConfigureInput struct with a mode
// field that is empty.
var ErrMissingMode = fastly.NewFieldError("Mode")

// ConfigureInput holds the details required by the API's
// 'UpdateConfiguration' operation.
type ConfigureInput struct {
	Mode string `json:"mode"`
}

// ConfigureOutput holds the details returned by the API from
// 'GetConfiguration' and 'UpdateConfiguration' operations.
type ConfigureOutput struct {
	products.ServiceConfigureOutput `mapstructure:",squash"`
	Configuration                   *configureOutputNested `mapstructure:"configuration"`
}

type configureOutputNested struct {
	Mode *string `mapstructure:"mode"`
}

// Get gets the status of the DDoS Protection product on the service.
func Get(ctx context.Context, c *fastly.Client, serviceID string) (EnableOutput, error) {
	return serviceproductcore.Get[EnableOutput](ctx, &serviceproductcore.GetInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Enable enables the DDoS Protection product on the service.
func Enable(ctx context.Context, c *fastly.Client, serviceID string, i EnableInput) (EnableOutput, error) {
	return serviceproductcore.Put[EnableOutput](ctx, &serviceproductcore.PutInput[EnableInput]{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
		Input:     i,
	})
}

// Disable disables the DDoS Protection product on the service.
func Disable(ctx context.Context, c *fastly.Client, serviceID string) error {
	return serviceproductcore.Delete(ctx, &serviceproductcore.DeleteInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// GetConfiguration gets the configuration of the DDoS Protection product on the service.
func GetConfiguration(ctx context.Context, c *fastly.Client, serviceID string) (ConfigureOutput, error) {
	return serviceproductcore.Get[ConfigureOutput](ctx, &serviceproductcore.GetInput{
		Client:        c,
		ProductID:     ProductID,
		ServiceID:     serviceID,
		URLComponents: []string{"configuration"},
	})
}

// UpdateConfiguration updates the configuration of the DDoS Protection product on the service.
func UpdateConfiguration(ctx context.Context, c *fastly.Client, serviceID string, i ConfigureInput) (ConfigureOutput, error) {
	if i.Mode == "" {
		return ConfigureOutput{}, ErrMissingMode
	}

	return serviceproductcore.Patch[ConfigureOutput](ctx, &serviceproductcore.PatchInput[ConfigureInput]{
		Client:        c,
		ProductID:     ProductID,
		ServiceID:     serviceID,
		URLComponents: []string{"configuration"},
		Input:         i,
	})
}

// NewEnableOutput is used to construct mock API output structures for
// use in tests.
func NewEnableOutput(serviceID string) EnableOutput {
	return products.NewServiceEnableOutput(ProductID, serviceID)
}
