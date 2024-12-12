package products

import "github.com/fastly/go-fastly/v9/fastly"

// GetInput specifies the information needed for the Get()
// function to perform the operation.
//
// Because Get operations produce output, this struct has a type
// parameter used to specify the type of the output; the parameter is
// not used within this structure, but specifying it at this level
// makes type inference simpler when the caller invokes the Get()
// function.
type GetInput[O ProductOutput] struct {
	Client        *fastly.Client
	ProductID     string
	ServiceID     string
	URLComponents []string
}

// Get implements a product-specific 'get' operation.
//
// This function requires the same type parameter as the GetInput
// struct, and that type is used to construct, populate, and return
// the output present in the response body.
func Get[O ProductOutput](i *GetInput[O]) (*O, error) {
	var err error
	if i.ServiceID == "" {
		return nil, fastly.ErrMissingServiceID
	}

	path := makeURL(i.ProductID, i.ServiceID, i.URLComponents)

	resp, err := i.Client.GetJSON(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *O
	if err = fastly.DecodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}

	return h, nil
}
