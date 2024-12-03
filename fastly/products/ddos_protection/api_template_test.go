package ddos_protection_test

import (
	"github.com/fastly/go-fastly/v9/fastly/products/ddos_protection"
)

type ConfigureInputTestCase struct {
	Name      string
	Input     ddos_protection.ConfigureInput
	WantError error
}

var ConfigureInputTestCases = map[string][]ConfigureInputTestCase{
	"valid": {
		{
			Name:  "valid",
			Input: ddos_protection.ConfigureInput{Mode: "off"},
		},
	},
	"invalid": {
		{
			Name:      "empty Mode",
			WantError: ddos_protection.ErrMissingMode,
		},
	},
}
