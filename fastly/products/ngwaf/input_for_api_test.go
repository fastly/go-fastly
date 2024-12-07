package ngwaf_test

import (
	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products/ngwaf"
	p "github.com/fastly/go-fastly/v9/internal/products"
)

var validEnable = []validateEnableInput{
	{
		name:  "valid",
		input: ngwaf.EnableInput{WorkspaceID: fastly.TestNGWAFWorkspaceID},
	},
}

var invalidEnable = []validateEnableInput{
	{
		name:      "empty WorkspaceID",
		wantError: ngwaf.ErrMissingWorkspaceID,
	},
}

var validUpdateConfiguration = []validateConfigureInput{
	{
		name:  "valid",
		input: ngwaf.ConfigureInput{},
	},
}

var invalidUpdateConfiguration = []validateConfigureInput{
	// there are no 'invalid' cases here, as all of the fields in
	// the ConfigureInput structure are optional
}

var serviceID = fastly.TestDeliveryServiceID

var functionalTests = []*fastly.FunctionalTestInput{
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase:         "before enablement",
		Executor:      ngwaf.Get,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
	p.TestEnable(&p.TestEnableInput[ngwaf.EnableInput, fastly.ProductEnablement]{
		ExecutorWithInput: ngwaf.Enable,
		Input:             &ngwaf.EnableInput{WorkspaceID: fastly.TestNGWAFWorkspaceID},
		ServiceID:         serviceID,
	}),
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase:     "after enablement",
		Executor:  ngwaf.Get,
		ServiceID: serviceID,
	}),
	p.TestUpdateConfiguration(&p.TestUpdateConfigurationInput[ngwaf.ConfigureInput, ngwaf.ConfigureOutput]{
		Executor:  ngwaf.UpdateConfiguration,
		Input:     &ngwaf.ConfigureInput{TrafficRamp: "45"},
		ServiceID: serviceID,
	}),
	p.TestGetConfiguration(&p.TestGetConfigurationInput[ngwaf.ConfigureOutput]{
		Executor:  ngwaf.GetConfiguration,
		ServiceID: serviceID,
	}),
	p.TestDisable(&p.TestDisableInput{
		Executor:  ngwaf.Disable,
		ServiceID: serviceID,
	}),
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase:         "after disablement",
		Executor:      ngwaf.Get,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
}
