package productcore

import (
	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products"
)

// GetInput specifies the information needed for the Get()
// function to perform the operation.
type GetInput struct {
	Client        *fastly.Client
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

	resp, err := i.Client.GetJSON(path, nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = fastly.DecodeBodyMap(resp.Body, &o)
	return
}
