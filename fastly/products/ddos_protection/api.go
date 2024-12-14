package ddos_protection

import (
	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products"
	"github.com/fastly/go-fastly/v9/internal/productcore"
)

const ProductID = "ddos_protection"

// ErrMissingMode is the error returned by the UpdateConfiguration
// function when it is passed a ConfigureInput struct with a mode
// field that is empty.
var ErrMissingMode = fastly.NewFieldError("Mode")

type ConfigureInput struct {
	Mode string `json:"mode"`
}

type ConfigureOutput struct {
	products.ConfigureOutput `mapstructure:",squash"`
	Configuration            *configureOutputNested `mapstructure:"configuration"`
}

type configureOutputNested struct {
	Mode *string `mapstructure:"mode"`
}

// Get gets the status of the DDoS Protection product on the service.
func Get(c *fastly.Client, serviceID string) (*products.EnableOutput, error) {
	return productcore.Get[*products.EnableOutput](&productcore.GetInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Enable enables the DDoS Protection product on the service.
func Enable(c *fastly.Client, serviceID string) (*products.EnableOutput, error) {
	return productcore.Put[*products.EnableOutput](&productcore.PutInput[*productcore.NullInput]{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Disable disables the DDoS Protection product on the service.
func Disable(c *fastly.Client, serviceID string) error {
	return productcore.Delete(&productcore.DeleteInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// GetConfiguration gets the configuration of the DDoS Protection product on the service.
func GetConfiguration(c *fastly.Client, serviceID string) (*ConfigureOutput, error) {
	return productcore.Get[*ConfigureOutput](&productcore.GetInput{
		Client:        c,
		ProductID:     ProductID,
		ServiceID:     serviceID,
		URLComponents: []string{"configuration"},
	})
}

// UpdateConfiguration updates the configuration of the DDoS Protection product on the service.
func UpdateConfiguration(c *fastly.Client, serviceID string, i *ConfigureInput) (*ConfigureOutput, error) {
	if i.Mode == "" {
		return nil, ErrMissingMode
	}

	return productcore.Patch[*ConfigureOutput](&productcore.PatchInput[*ConfigureInput]{
		Client:        c,
		ProductID:     ProductID,
		ServiceID:     serviceID,
		URLComponents: []string{"configuration"},
		Input:         i,
	})
}
