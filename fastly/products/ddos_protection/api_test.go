package ddos_protection_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products/ddos_protection"
	p "github.com/fastly/go-fastly/v9/internal/products"
)

var serviceID = fastly.TestDeliveryServiceID

var functionalTests = []*fastly.FunctionalTestInput{
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase:         "before enablement",
		Executor:      ddos_protection.Get,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
	p.TestEnable(&p.TestEnableInput[p.NullInput, fastly.ProductEnablement]{
		ExecutorNoInput: ddos_protection.Enable,
		ServiceID:       serviceID,
	}),
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase:     "after enablement",
		Executor:  ddos_protection.Get,
		ServiceID: serviceID,
	}),
	p.TestUpdateConfiguration(&p.TestUpdateConfigurationInput[ddos_protection.ConfigureInput, ddos_protection.ConfigureOutput]{
		Executor:  ddos_protection.UpdateConfiguration,
		Input:     &ddos_protection.ConfigureInput{Mode: "logging"},
		ServiceID: serviceID,
	}),
	p.TestGetConfiguration(&p.TestGetConfigurationInput[ddos_protection.ConfigureOutput]{
		Executor:  ddos_protection.GetConfiguration,
		ServiceID: serviceID,
	}),
	p.TestDisable(&p.TestDisableInput{
		Executor:  ddos_protection.Disable,
		ServiceID: serviceID,
	}),
	p.TestGet(&p.TestGetInput[fastly.ProductEnablement]{
		Phase:         "after disablement",
		Executor:      ddos_protection.Get,
		ServiceID:     serviceID,
		ExpectFailure: true,
	}),
}

func TestEnablementAndConfiguration(t *testing.T) {
	fastly.ExecuteFunctionalTests(t, functionalTests)
}
