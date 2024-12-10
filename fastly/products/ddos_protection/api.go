package ddos_protection

import (
	"github.com/fastly/go-fastly/v9/fastly"
	// fp is 'fastly products' package
	fp "github.com/fastly/go-fastly/v9/fastly/products"
	// ip is 'internal products' package
	ip "github.com/fastly/go-fastly/v9/internal/products"
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
	fp.ConfigureOutput `mapstructure:",squash"`
	Configuration      *configureOutputNested `mapstructure:"configuration"`
}

type configureOutputNested struct {
	Mode *string `mapstructure:"mode"`
}

// Get gets the status of the DDoS Protection product on the service.
func Get(c *fastly.Client, serviceID string) (*fp.EnableOutput, error) {
	return ip.Get(&ip.GetInput[fp.EnableOutput]{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Enable enables the DDoS Protection product on the service.
func Enable(c *fastly.Client, serviceID string) (*fp.EnableOutput, error) {
	return ip.Put(&ip.PutInput[ip.NullInput, fp.EnableOutput]{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Disable disables the DDoS Protection product on the service.
func Disable(c *fastly.Client, serviceID string) error {
	return ip.Delete(&ip.DeleteInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// GetConfiguration gets the configuration of the DDoS Protection product on the service.
func GetConfiguration(c *fastly.Client, serviceID string) (*ConfigureOutput, error) {
	return ip.Get(&ip.GetInput[ConfigureOutput]{
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

	return ip.Patch(&ip.PatchInput[ConfigureInput, ConfigureOutput]{
		Client:        c,
		ProductID:     ProductID,
		ServiceID:     serviceID,
		URLComponents: []string{"configuration"},
		Input:         i,
	})
}
