package productcore

import (
	"context"

	"github.com/fastly/go-fastly/v19/fastly"
)

// DeleteInput specifies the information needed for the Delete
// function to perform the operation.
type DeleteInput struct {
	Client        *fastly.Client
	ProductID     string
	URLComponents []string
}

// Delete implements a product-specific 'delete' operation. Since this
// operation does not accept any input or produce any output (other
// than a potential error), this function does not have any type
// parameters.
func Delete(ctx context.Context, i *DeleteInput) error {
	path := makeURL(i.ProductID, i.URLComponents)

	resp, err := i.Client.Delete(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
