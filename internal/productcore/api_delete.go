package productcore

import (
	"context"

	"github.com/fastly/go-fastly/v10/fastly"
)

// DeleteInput specifies the information needed for the Delete
// function to perform the operation.
type DeleteInput struct {
	Client *fastly.Client
	// Context, if supplied, will be used as the Request's context.
	Context       *context.Context
	ProductID     string
	ServiceID     string
	URLComponents []string
}

// Delete implements a product-specific 'delete' operation. Since this
// operation does not accept any input or produce any output (other
// than a potential error), this function does not have any type
// parameters.
func Delete(i *DeleteInput) error {
	if i.ServiceID == "" {
		return fastly.ErrMissingServiceID
	}

	path := makeURL(i.ProductID, i.ServiceID, i.URLComponents)

	resp, err := i.Client.Delete(path, fastly.CreateRequestOptions(i.Context))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
