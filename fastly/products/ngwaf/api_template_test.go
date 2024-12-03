package ngwaf_test

import (
	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products/ngwaf"
)

type EnableInputTestCase struct {
	Name      string
	Input     ngwaf.EnableInput
	WantError error
}

var EnableInputTestCases = map[string][]EnableInputTestCase{
	"valid": {
		{
			Name:  "valid",
			Input: ngwaf.EnableInput{WorkspaceID: fastly.TestNGWAFWorkspaceID},
		},
	},
	"invalid": {
		{
			Name:      "empty WorkspaceID",
			WantError: ngwaf.ErrMissingWorkspaceID,
		},
	},
}

type ConfigureInputTestCase struct {
	Name      string
	Input     ngwaf.ConfigureInput
	WantError error
}

var ConfigureInputTestCases = map[string][]ConfigureInputTestCase{
	"valid": {
		{
			Name:  "valid",
			Input: ngwaf.ConfigureInput{},
		},
	},
	// there are no 'invalid' cases here, as all of the fields in
	// the ConfigureInput structure are optional
}
