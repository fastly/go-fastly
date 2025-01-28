package domaininspector

import (
	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products"
	"github.com/fastly/go-fastly/v9/internal/productcore"
)

const (
	ProductID   = "domain_inspector"
	ProductName = "Domain Inspector"
)

// EnableOutput holds the details returned by the API from 'Get' and
// 'Enable' operations; this alias exists to ensure that users of this
// package will have a stable name to reference.
type EnableOutput = products.EnableOutput

// Get gets the status of the Domain Inspector product on the service.
func Get(c *fastly.Client, serviceID string) (EnableOutput, error) {
	return productcore.Get[EnableOutput](&productcore.GetInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Enable enables the Domain Inspector product on the service.
func Enable(c *fastly.Client, serviceID string) (EnableOutput, error) {
	return productcore.Put[EnableOutput](&productcore.PutInput[products.NullInput]{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Disable disables the Domain Inspector product on the service.
func Disable(c *fastly.Client, serviceID string) error {
	return productcore.Delete(&productcore.DeleteInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// NewEnableOutput is used to construct mock API output structures for
// use in tests.
func NewEnableOutput(serviceID string) EnableOutput {
	return products.NewEnableOutput(ProductID, serviceID)
}
