package apidiscovery

import (
	"context"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/internal/serviceproductcore"
)

const (
	ProductID   = "api_discovery"
	ProductName = "API Discovery"
)

// EnableOutput holds the details returned by the API from 'Get' and
// 'Enable' operations; this alias exists to ensure that users of this
// package will have a stable name to reference.
type EnableOutput = products.ServiceEnableOutput

// Get gets the status of the API Discovery product on the service.
func Get(ctx context.Context, c *fastly.Client, serviceID string) (EnableOutput, error) {
	return serviceproductcore.Get[EnableOutput](ctx, &serviceproductcore.GetInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Enable enables the API Discovery product on the service.
func Enable(ctx context.Context, c *fastly.Client, serviceID string) (EnableOutput, error) {
	return serviceproductcore.Put[EnableOutput](ctx, &serviceproductcore.PutInput[products.ServiceNullInput]{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Disable disables the API Discovery product on the service.
func Disable(ctx context.Context, c *fastly.Client, serviceID string) error {
	return serviceproductcore.Delete(ctx, &serviceproductcore.DeleteInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// NewEnableOutput is used to construct mock API output structures for
// use in tests.
func NewEnableOutput(serviceID string) EnableOutput {
	return products.NewServiceEnableOutput(ProductID, serviceID)
}
