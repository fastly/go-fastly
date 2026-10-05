package contacts

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/jsonapi"

	"github.com/fastly/go-fastly/v15/fastly"
)

// ListInput specifies the information needed for the List() function to
// perform the operation.
type ListInput struct {
	// CustomerID is the alphanumeric identifier of the customer (required).
	CustomerID *string
}

// List retrieves all contacts for the given customer.
func List(ctx context.Context, c *fastly.Client, i *ListInput) ([]*Contact, error) {
	if i.CustomerID == nil {
		return nil, fastly.ErrMissingCustomerID
	}

	path := fastly.ToSafeURL("customer", *i.CustomerID, "contacts")

	opts := fastly.CreateRequestOptions()
	opts.Headers["Accept"] = jsonapi.MediaType

	resp, err := c.Get(ctx, path, opts)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Contact)))
	if err != nil {
		return nil, err
	}

	cs := make([]*Contact, len(data))
	for idx := range data {
		typed, ok := data[idx].(*Contact)
		if !ok {
			return nil, fmt.Errorf("unexpected response type: %T", data[idx])
		}
		cs[idx] = typed
	}

	return cs, nil
}
