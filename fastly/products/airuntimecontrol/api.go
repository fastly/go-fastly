package airuntimecontrol

import (
	"context"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/internal/productcore"
)

const (
	ProductID   = "ai_runtime_control"
	ProductName = "AI Runtime Control"
)

// EnableOutput holds the details returned by the API from 'Get' and
// 'Enable' operations; this alias exists to ensure that users of this
// package will have a stable name to reference.
type EnableOutput = products.EnableOutput

// Get gets the status of the AI Runtime Control product on the account.
func Get(ctx context.Context, c *fastly.Client) (EnableOutput, error) {
	return productcore.Get[EnableOutput](ctx, &productcore.GetInput{
		Client:    c,
		ProductID: ProductID,
	})
}

// Enable enables the AI Runtime Control product on the account.
func Enable(ctx context.Context, c *fastly.Client) (EnableOutput, error) {
	return productcore.Put[EnableOutput](ctx, &productcore.PutInput[products.NullInput]{
		Client:    c,
		ProductID: ProductID,
	})
}

// Disable disables the AI Runtime Control product on the account.
func Disable(ctx context.Context, c *fastly.Client) error {
	return productcore.Delete(ctx, &productcore.DeleteInput{
		Client:    c,
		ProductID: ProductID,
	})
}

// NewEnableOutput is used to construct mock API output structures for
// use in tests.
func NewEnableOutput() EnableOutput {
	return products.NewEnableOutput(ProductID)
}
