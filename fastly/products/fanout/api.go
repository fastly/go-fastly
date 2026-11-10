package fanout

import (
	"context"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/internal/serviceproductcore"
)

const (
	ProductID   = "fanout"
	ProductName = "Fanout"
)

// ServiceEnableOutput holds the details returned by the API from 'Get' and
// 'Enable' operations; this alias exists to ensure that users of this
// package will have a stable name to reference.
type ServiceEnableOutput = products.ServiceEnableOutput

// Get gets the status of the Fanout product on the service.
func Get(ctx context.Context, c *fastly.Client, serviceID string) (ServiceEnableOutput, error) {
	return serviceproductcore.Get[ServiceEnableOutput](ctx, &serviceproductcore.GetInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Enable enables the Fanout product on the service.
func Enable(ctx context.Context, c *fastly.Client, serviceID string) (ServiceEnableOutput, error) {
	return serviceproductcore.Put[ServiceEnableOutput](ctx, &serviceproductcore.PutInput[products.ServiceNullInput]{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Disable disables the Fanout product on the service.
func Disable(ctx context.Context, c *fastly.Client, serviceID string) error {
	return serviceproductcore.Delete(ctx, &serviceproductcore.DeleteInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// NewServiceEnableOutput is used to construct mock API output structures for
// use in tests.
func NewServiceEnableOutput(serviceID string) ServiceEnableOutput {
	return products.NewServiceEnableOutput(ProductID, serviceID)
}
