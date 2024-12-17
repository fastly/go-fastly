package log_explorer_insights

import (
	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products"
	"github.com/fastly/go-fastly/v9/internal/productcore"
)

const ProductID = "log_explorer_insights"

// Get gets the status of the Log Explorer Bot Management Insights product on the service.
func Get(c *fastly.Client, serviceID string) (*products.EnableOutput, error) {
	return productcore.Get[*products.EnableOutput](&productcore.GetInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Enable enables the Log Explorer Bot Management Insights product on the service.
func Enable(c *fastly.Client, serviceID string) (*products.EnableOutput, error) {
	return productcore.Put[*products.EnableOutput](&productcore.PutInput[*productcore.NullInput]{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Disable disables the Log Explorer Bot Management Insights product on the service.
func Disable(c *fastly.Client, serviceID string) error {
	return productcore.Delete(&productcore.DeleteInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}
