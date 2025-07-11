package brotlicompression

import (
	"context"

	"github.com/fastly/go-fastly/v10/fastly"
	"github.com/fastly/go-fastly/v10/fastly/products"
	"github.com/fastly/go-fastly/v10/internal/productcore"
)

const (
	ProductID   = "brotli_compression"
	ProductName = "Brotli Compression"
)

// EnableOutput holds the details returned by the API from 'Get' and
// 'Enable' operations; this alias exists to ensure that users of this
// package will have a stable name to reference.
type EnableOutput = products.EnableOutput

// Get gets the status of the Brotli Compression product on the service.
func Get(ctx context.Context, c *fastly.Client, serviceID string) (EnableOutput, error) {
	return productcore.Get[EnableOutput](ctx, &productcore.GetInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Enable enables the Brotli Compression product on the service.
func Enable(ctx context.Context, c *fastly.Client, serviceID string) (EnableOutput, error) {
	return productcore.Put[EnableOutput](ctx, &productcore.PutInput[products.NullInput]{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Disable disables the Brotli Compression product on the service.
func Disable(ctx context.Context, c *fastly.Client, serviceID string) error {
	return productcore.Delete(ctx, &productcore.DeleteInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// NewEnableOutput is used to construct mock API output structures for
// use in tests.
func NewEnableOutput(serviceID string) EnableOutput {
	return products.NewEnableOutput(ProductID, serviceID)
}
