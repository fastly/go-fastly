//go:generate rm -f api.go
//go:generate service_linked_product -api

package ddos_protection

import (
	"github.com/fastly/go-fastly/v9/fastly"
)

const (
	ProductName = "DDoS Protection"
	ProductID   = "ddos_protection"
)

// ErrMissingMode is the error returned by the UpdateConfiguration
// function when it is passed a ConfigureInput struct with a mode
// field that is empty.
var ErrMissingMode = fastly.NewFieldError("Mode")

type ConfigureInput struct {
	Mode string `json:"mode"`
}

func (i *ConfigureInput) Validate() error {
	if i.Mode == "" {
		return ErrMissingMode
	}
	return nil
}

type ConfigureOutput struct {
	fastly.ProductConfiguration
	Configuration *configureOutputNested `mapstructure:"configuration"`
}

type configureOutputNested struct {
	Mode *string `mapstructure:"mode"`
}
