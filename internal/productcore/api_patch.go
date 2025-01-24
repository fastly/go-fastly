package productcore

import (
	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products"
)

// PatchInput specifies the information needed for the Patch()
// function to perform the operation.
//
// Because Patch operations accept input, this struct has a type
// parameter used to specify the type of the input structure.
type PatchInput[I any] struct {
	Client        *fastly.Client
	ProductID     string
	ServiceID     string
	URLComponents []string
	Input         I
}

// Patch implements a product-specific 'patch' operation.
//
// This function requires the same type parameter as the PatchInput
// struct; the input type parameter is used to marshal the input data
// into the request body (encoded as JSON).
//
// It also requires a type parameter which is a pointer to an
// struct which matches the ProductOutput interface, and that type
// is used to construct, populate, and return the output present in
// the response body.
func Patch[O products.ProductOutput, I any](i *PatchInput[I]) (o O, err error) {
	if i.ServiceID == "" {
		err = fastly.ErrMissingServiceID
		return
	}

	path := makeURL(i.ProductID, i.ServiceID, i.URLComponents)

	resp, err := i.Client.PatchJSON(path, i.Input, nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = fastly.DecodeBodyMap(resp.Body, &o)
	return
}
