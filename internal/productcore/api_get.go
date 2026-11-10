package productcore

import (
	"context"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products"
)

// GetInput specifies the information needed for the Get()
// function to perform the operation.
type GetInput struct {
	Client        *fastly.Client
	ProductID     string
	URLComponents []string
}

// Get implements a product-specific 'get' operation.
//
// This function requires a type parameter which is a pointer to an
// struct which matches the ProductOutput interface, and that type
// is used to construct, populate, and return the output present in
// the response body.
func Get[O products.ProductOutput](ctx context.Context, i *GetInput) (o O, err error) {
	path := makeURL(i.ProductID, i.URLComponents)

	resp, err := i.Client.GetJSON(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = fastly.DecodeBodyMap(resp.Body, &o)
	return
}
