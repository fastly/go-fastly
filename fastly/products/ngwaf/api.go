package ngwaf

import (
	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products"
	"github.com/fastly/go-fastly/v9/internal/productcore"
)

const (
	ProductID   = "ngwaf"
	ProductName = "Next-Gen WAF"
)

// EnableInput holds the details required by the API's 'Enable'
// operation.
type EnableInput struct {
	WorkspaceID string `json:"workspace_id"`
}

// EnableOutput holds the details returned by the API from 'Get' and
// 'Enable' operations; this alias exists to ensure that users of this
// package will have a stable name to reference.
type EnableOutput = products.EnableOutput

// ErrMissingWorkspaceID is the error returned by the Enable function
// when it is passed an EnableInput struct with a WorkspaceID field
// that is empty.
var ErrMissingWorkspaceID = fastly.NewFieldError("WorkspaceID")

// ConfigureInput holds the details required by the API's
// 'UpdateConfiguration' operation.
type ConfigureInput struct {
	WorkspaceID string `json:"workspace_id,omitempty"`
	TrafficRamp string `json:"traffic_ramp,omitempty"`
}

// ConfigureOutput holds the details returned by the API from
// 'GetConfiguration' and 'UpdateConfiguration' operations.
type ConfigureOutput struct {
	products.ConfigureOutput `mapstructure:",squash"`
	Configuration            *configureOutputNested `mapstructure:"configuration"`
}

type configureOutputNested struct {
	WorkspaceID *string `mapstructure:"workspace_id,omitempty"`
	TrafficRamp *string `mapstructure:"traffic_ramp,omitempty"`
}

// Get gets the status of the Next-Gen WAF product on the service.
func Get(c *fastly.Client, serviceID string) (EnableOutput, error) {
	return productcore.Get[EnableOutput](&productcore.GetInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// Enable enables the Next-Gen WAF product on the service.
func Enable(c *fastly.Client, serviceID string, i EnableInput) (EnableOutput, error) {
	if i.WorkspaceID == "" {
		return EnableOutput{}, ErrMissingWorkspaceID
	}

	return productcore.Put[EnableOutput](&productcore.PutInput[EnableInput]{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
		Input:     i,
	})
}

// Disable disables the Next-Gen WAF product on the service.
func Disable(c *fastly.Client, serviceID string) error {
	return productcore.Delete(&productcore.DeleteInput{
		Client:    c,
		ProductID: ProductID,
		ServiceID: serviceID,
	})
}

// GetConfiguration gets the configuration of the Next-Gen WAF product on the service.
func GetConfiguration(c *fastly.Client, serviceID string) (ConfigureOutput, error) {
	return productcore.Get[ConfigureOutput](&productcore.GetInput{
		Client:        c,
		ProductID:     ProductID,
		ServiceID:     serviceID,
		URLComponents: []string{"configuration"},
	})
}

// UpdateConfiguration updates the configuration of the Next-Gen WAF product on the service.
func UpdateConfiguration(c *fastly.Client, serviceID string, i ConfigureInput) (ConfigureOutput, error) {
	return productcore.Patch[ConfigureOutput](&productcore.PatchInput[ConfigureInput]{
		Client:        c,
		ProductID:     ProductID,
		ServiceID:     serviceID,
		URLComponents: []string{"configuration"},
		Input:         i,
	})
}
