package domain_inspector

import (
	"github.com/fastly/go-fastly/v9/fastly"
	// fp is 'fastly products' package
	fp "github.com/fastly/go-fastly/v9/fastly/products"
	// ip is 'internal products' package
	ip "github.com/fastly/go-fastly/v9/internal/products"
)

const ProductID = "domain_inspector"

// Get gets the status of the Domain Inspector product on the service.
func Get(c *fastly.Client, serviceID string) (*fp.EnableOutput, error) {
	return ip.Get(&ip.GetInput[fp.EnableOutput]{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Enable enables the Domain Inspector product on the service.
func Enable(c *fastly.Client, serviceID string) (*fp.EnableOutput, error) {
	return ip.Put(&ip.PutInput[ip.NullInput, fp.EnableOutput]{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Disable disables the Domain Inspector product on the service.
func Disable(c *fastly.Client, serviceID string) error {
	return ip.Delete(&ip.DeleteInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}
