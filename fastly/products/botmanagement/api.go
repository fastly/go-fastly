package botmanagement

import (
	"context"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/internal/serviceproductcore"
)

const (
	ProductID   = "bot_management"
	ProductName = "Bot Management"
)

// ServiceEnableOutput holds the details returned by the API from 'Get' and
// 'Enable' operations; this alias exists to ensure that users of this
// package will have a stable name to reference.
type ServiceEnableOutput = products.ServiceEnableOutput

// ErrMissingContentGuard is the error returned by the UpdateConfiguration
// function when it is passed a ConfigureInput struct with a contentguard
// field that is empty.
var ErrMissingContentGuard = fastly.NewFieldError("ContentGuard")

// ConfigureInput holds the details required by the API's
// 'UpdateConfiguration' operation.
type ConfigureInput struct {
	ContentGuard string `json:"contentguard"`
}

// ServiceConfigureOutput holds the details returned by the API from
// 'GetConfiguration' and 'UpdateConfiguration' operations.
type ServiceConfigureOutput struct {
	products.ServiceConfigureOutput `mapstructure:",squash"`
	Configuration                   *configureOutputNested `mapstructure:"configuration"`
}

type configureOutputNested struct {
	ContentGuard *string `mapstructure:"contentguard"`
}

// Get gets the status of the Bot Management product on the service.
func Get(ctx context.Context, c *fastly.Client, serviceID string) (ServiceEnableOutput, error) {
	return serviceproductcore.Get[ServiceEnableOutput](ctx, &serviceproductcore.GetInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Enable enables the Bot Management product on the service.
func Enable(ctx context.Context, c *fastly.Client, serviceID string) (ServiceEnableOutput, error) {
	return serviceproductcore.Put[ServiceEnableOutput](ctx, &serviceproductcore.PutInput[products.ServiceNullInput]{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Disable disables the Bot Management product on the service.
func Disable(ctx context.Context, c *fastly.Client, serviceID string) error {
	return serviceproductcore.Delete(ctx, &serviceproductcore.DeleteInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// GetConfiguration gets the configuration of the Bot Management product on the service.
func GetConfiguration(ctx context.Context, c *fastly.Client, serviceID string) (ServiceConfigureOutput, error) {
	return serviceproductcore.Get[ServiceConfigureOutput](ctx, &serviceproductcore.GetInput{
		Client:        c,
		ProductID:     ProductID,
		ServiceID:     serviceID,
		URLComponents: []string{"configuration"},
	})
}

// UpdateConfiguration updates the configuration of the Bot Management product on the service.
func UpdateConfiguration(ctx context.Context, c *fastly.Client, serviceID string, i ConfigureInput) (ServiceConfigureOutput, error) {
	if i.ContentGuard == "" {
		return ServiceConfigureOutput{}, ErrMissingContentGuard
	}

	return serviceproductcore.Patch[ServiceConfigureOutput](ctx, &serviceproductcore.PatchInput[ConfigureInput]{
		Client:        c,
		ProductID:     ProductID,
		ServiceID:     serviceID,
		URLComponents: []string{"configuration"},
		Input:         i,
	})
}

// NewServiceEnableOutput is used to construct mock API output structures for
// use in tests.
func NewServiceEnableOutput(serviceID string) ServiceEnableOutput {
	return products.NewServiceEnableOutput(ProductID, serviceID)
}
