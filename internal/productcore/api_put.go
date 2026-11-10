package productcore

import (
	"context"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products"
)

// PutInput specifies the information needed for the Put()
// function to perform the operation.
//
// Because Put operations accept input, this struct has a type
// parameter used to specify the type of the input structure.
type PutInput[I any] struct {
	Client        *fastly.Client
	ProductID     string
	URLComponents []string
	Input         I
}

// Put implements a product-specific 'put' operation.
//
// This function requires the same type parameter as the PutInput
// struct; the input type parameter is used to marshal the input data
// into the request body (encoded as JSON).
//
// It also requires a type parameter which is a pointer to an
// struct which matches the ProductOutput interface, and that type
// is used to construct, populate, and return the output present in
// the response body.
func Put[O products.ProductOutput, I any](ctx context.Context, i *PutInput[I]) (o O, err error) {
	path := makeURL(i.ProductID, i.URLComponents)

	resp, err := i.Client.PutJSON(ctx, path, i.Input, fastly.CreateRequestOptions())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = fastly.DecodeBodyMap(resp.Body, &o)
	return
}
