package domain_inspector

import (
	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/internal/productcore"
)

const (
	ProductID   = "domain_inspector"
	ProductName = "Domain Inspector"
)

type EnableOutput = productcore.EnableOutput

// Get gets the status of the Domain Inspector product on the service.
func Get(c *fastly.Client, serviceID string) (*EnableOutput, error) {
	return productcore.Get[*EnableOutput](&productcore.GetInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Enable enables the Domain Inspector product on the service.
func Enable(c *fastly.Client, serviceID string) (*EnableOutput, error) {
	return productcore.Put[*EnableOutput](&productcore.PutInput[*productcore.NullInput]{
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
