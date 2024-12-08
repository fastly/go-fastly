package products

import "github.com/fastly/go-fastly/v9/fastly"

// PatchInput specifies the information needed for the Patch()
// function to perform the operation.
//
// Because Patch operations accept input and produce output, this struct
// has two type parameters used to specify the types of the input and
// output; the output parameter is not used within this structure, but
// specifying it at this level makes type inference simpler when the
// caller invokes the Patch() function.
type PatchInput[I, O any] struct {
	Client        *fastly.Client
	ProductID     string
	ServiceID     string
	URLComponents []string
	Input         *I
}

// Patch implements a product-specific 'patch' operation.
//
// This function requires the same type parameters as the PatchInput
// struct; the input type parameter is used to marshal the input data
// into the request body (encoded as JSON), and the output type
// parameter is used to construct, populate, and return the output
// present in the response body.
func Patch[I, O any](i *PatchInput[I, O]) (*O, error) {
	var err error
	if i.ServiceID == "" {
		return nil, fastly.ErrMissingServiceID
	}

	path := makeURL(i.ProductID, i.ServiceID, i.URLComponents)

	resp, err := i.Client.PatchJSON(path, i.Input, nil)
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
