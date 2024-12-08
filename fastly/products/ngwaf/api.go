package ngwaf

import (
	"github.com/fastly/go-fastly/v9/fastly"
	// fp is 'fastly products' package
	fp "github.com/fastly/go-fastly/v9/fastly/products"
	// ip is 'internal products' package
	ip "github.com/fastly/go-fastly/v9/internal/products"
)

const ProductID = "ngwaf"

// ErrMissingWorkspaceID is the error returned by the Enable function
// when it is passed an EnableInput struct with a WorkspaceID field
// that is empty.
var ErrMissingWorkspaceID = fastly.NewFieldError("WorkspaceID")

type EnableInput struct {
	WorkspaceID string `json:"workspace_id"`
}

type ConfigureInput struct {
	WorkspaceID string `json:"workspace_id,omitempty"`
	TrafficRamp string `json:"traffic_ramp,omitempty"`
}

type ConfigureOutput struct {
	fp.ConfigureOutput `mapstructure:",squash"`
	Configuration      *configureOutputNested `mapstructure:"configuration"`
}

type configureOutputNested struct {
	WorkspaceID *string `mapstructure:"workspace_id,omitempty"`
	TrafficRamp *string `mapstructure:"traffic_ramp,omitempty"`
}

// Get gets the status of the Next-Gen WAF product on the service.
func Get(c *fastly.Client, serviceID string) (*fp.EnableOutput, error) {
	return ip.Get(&ip.GetInput[fp.EnableOutput]{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Enable enables the Next-Gen WAF product on the service.
func Enable(c *fastly.Client, serviceID string, i *EnableInput) (*fp.EnableOutput, error) {
	if i.WorkspaceID == "" {
		return nil, ErrMissingWorkspaceID
	}

	return ip.Put(&ip.PutInput[EnableInput, fp.EnableOutput]{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
		Input:     i,
	})
}

// Disable disables the Next-Gen WAF product on the service.
func Disable(c *fastly.Client, serviceID string) error {
	return ip.Delete(&ip.DeleteInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// GetConfiguration gets the configuration of the Next-Gen WAF product on the service.
func GetConfiguration(c *fastly.Client, serviceID string) (*ConfigureOutput, error) {
	return ip.Get(&ip.GetInput[ConfigureOutput]{
		Client:        c,
		ProductID:     ProductID,
		ServiceID:     serviceID,
		URLComponents: []string{"configuration"},
	})
}

// UpdateConfiguration updates the configuration of the Next-Gen WAF product on the service.
func UpdateConfiguration(c *fastly.Client, serviceID string, i *ConfigureInput) (*ConfigureOutput, error) {
	return ip.Patch(&ip.PatchInput[ConfigureInput, ConfigureOutput]{
		Client:        c,
		ProductID:     ProductID,
		ServiceID:     serviceID,
		URLComponents: []string{"configuration"},
		Input:         i,
	})
}
