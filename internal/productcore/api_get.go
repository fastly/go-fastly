package productcore

import (
	"context"

	"github.com/fastly/go-fastly/v10/fastly"
	"github.com/fastly/go-fastly/v10/fastly/products"
)

// GetInput specifies the information needed for the Get()
// function to perform the operation.
type GetInput struct {
	Client *fastly.Client
	// Context, if supplied, will be used as the Request's context.
	Context       *context.Context
	ProductID     string
	ServiceID     string
	URLComponents []string
}

// Get implements a product-specific 'get' operation.
//
// This function requires a type parameter which is a pointer to an
// struct which matches the ProductOutput interface, and that type
// is used to construct, populate, and return the output present in
// the response body.
func Get[O products.ProductOutput](i *GetInput) (o O, err error) {
	if i.ServiceID == "" {
		err = fastly.ErrMissingServiceID
		return
	}

	path := makeURL(i.ProductID, i.ServiceID, i.URLComponents)

	resp, err := i.Client.GetJSON(path, fastly.CreateRequestOptions(i.Context))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = fastly.DecodeBodyMap(resp.Body, &o)
	return
}
