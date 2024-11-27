//go:generate rm -f api.go api_test.go
//go:generate service_linked_product

package ddos_protection

import (
	"github.com/fastly/go-fastly/v9/fastly"
)

const (
	ProductName = "DDOS Protection"
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
	Configuration *ConfigureOutputNested `mapstructure:"configuration"`
}

type ConfigureOutputNested struct {
	Mode *string `mapstructure:"mode"`
}

type ConfigureInputTestCase struct {
	Name      string
	Input     ConfigureInput
	WantError error
}

var ConfigureInputTestCases = map[string][]ConfigureInputTestCase{
	"valid": {
		{
			Name:  "valid",
			Input: ConfigureInput{Mode: "off"},
		},
	},
	"invalid": {
		{
			Name:      "empty Mode",
			WantError: ErrMissingMode,
		},
	},
}
